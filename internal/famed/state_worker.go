package famed

import (
	"context"
	"github.com/phuslu/log"

	"github.com/morphysm/famed-github-backend/internal/config"
	model2 "github.com/morphysm/famed-github-backend/internal/repositories/github/model"
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
			famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
			issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, installation.Account.Login, repoName, []string{famedLabel.Name}, nil)
			if err != nil {
				log.Error().Err(err).Msgf("[CleanState] error while fetching issues for %s/%s", installation.Account.Login, repoName)
			}

			go func(owner string, repoName string, issues []model2.Issue) {
				gH.updateRewardComments(ctx, owner, repoName, issues, nil)
			}(installation.Account.Login, repoName, issues)
			go func(owner string, repoName string, issues []model2.Issue) {
				gH.updateEligibleComments(ctx, owner, repoName, issues, nil)
			}(installation.Account.Login, repoName, issues)
		}
	}
}
