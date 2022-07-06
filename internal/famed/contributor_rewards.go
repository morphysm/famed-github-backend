package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/config"
	famedModel "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
)

// GetContributorRewards returns a list of rewards for a given contributor.
func (gH *githubHandler) GetContributorRewards(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		contributor = c.Param("contributor")
	)

	installations, err := gH.githubAppClient.GetInstallations(ctx)
	if err != nil {
		return err
	}

	response := EligibleIssuesResponse{Repos: []Repo{}}
	for _, installation := range installations {
		owner := installation.Account.Login
		// TODO document drawback of contributor having to be assigned to issue
		repos, err := gH.githubInstallationClient.GetReposByOwner(ctx, owner)
		if err != nil {
			return err
		}

		issueState := model.Closed
		issueOptions := providers.IssueListByRepoOptions{
			Labels:   []string{gH.famedConfig.Labels[config.FamedLabelKey].Name},
			State:    &issueState,
			Assignee: &contributor,
		}
		for _, repo := range repos {
			repoResp := Repo{Name: repo, EligibleIssues: []EligibleIssue{}}
			issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repo, issueOptions)
			if err != nil {
				return err
			}

			for _, issue := range issues {
				enrichedIssue := gH.githubInstallationClient.EnrichIssue(ctx, owner, repo, issue)

				// TODO compress into githubHandler struct
				rewardStructure := famedModel.NewRewardStructure(gH.famedConfig.Rewards, gH.famedConfig.DaysToFix, 2)
				boardOptions := famedModel.NewBoardOptions(gH.famedConfig.Currency, rewardStructure, gH.now())
				contributors, err := famedModel.NewBlueTeamFromIssue(enrichedIssue, boardOptions)

				if len(contributors) != 0 && err == nil {
					repoResp.EligibleIssues = append(repoResp.EligibleIssues, EligibleIssue{ID: issue.ID, Number: issue.Number, HTMLURL: issue.HTMLURL, Contributors: contributors})
				}
			}

			response.Repos = append(response.Repos, repoResp)
		}
	}

	return c.JSON(http.StatusOK, response)
}
