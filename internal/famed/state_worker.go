package famed

import (
	"context"
	"github.com/phuslu/log"

	famedModel "github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// CleanState iterates over all issues and updates their comments if necessary.
func (gH *githubHandler) CleanState() {
	log.Info().Msgf("[CleanState] running clean up...")

	ctx := context.Background()
	installations, err := gH.githubAppClient.GetInstallations(ctx)
	if err != nil {
		log.Error().Err(err).Msg("[CleanState] error while getting installations")
	}

	for _, installation := range installations {
		// Check if installation client is set up and if necessary add client
		if !gH.githubInstallationClient.CheckInstallation(installation.Account.Login) {
			err := gH.githubInstallationClient.AddInstallation(installation.Account.Login, installation.ID)
			if err != nil {
				log.Error().Err(err).Msg("[CleanState] error while adding github")
				continue
			}
		}

		repos, err := gH.githubInstallationClient.GetRepos(ctx, installation.Account.Login)
		if err != nil {
			log.Error().Err(err).Msg("[CleanState] error while getting repos")
			continue
		}

		for _, repoName := range repos {
			issues, err := gH.githubInstallationClient.GetEnrichedIssues(ctx, installation.Account.Login, repoName, famedModel.Opened)
			if err != nil {
				log.Error().Err(err).Msgf("[CleanState] error while fetching issues for %s/%s", installation.Account.Login, repoName)
			}

			commentsIssues := make(map[*famedModel.EnrichedIssue][]famedModel.IssueComment, len(issues))
			for _, issue := range issues {
				comments, _ := gH.githubInstallationClient.GetComments(ctx, installation.Account.Login, repoName, issue.Number)
				commentsIssues[&issue] = comments
			}

			go func(owner string, repoName string, issues map[*famedModel.EnrichedIssue][]famedModel.IssueComment) {
				gH.updateRewardComments(ctx, owner, repoName, commentsIssues, nil)
			}(installation.Account.Login, repoName, commentsIssues)
			go func(owner string, repoName string, issues map[*famedModel.EnrichedIssue][]famedModel.IssueComment) {
				gH.updateEligibleComments(ctx, owner, repoName, commentsIssues, nil)
			}(installation.Account.Login, repoName, commentsIssues)
		}
	}
}
