package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/config"
	famedModel "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
)

type EligibleIssuesResponse struct {
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Name           string          `json:"name"`
	EligibleIssues []EligibleIssue `json:"eligibleIssues"`
}

// EligibleIssue represents an issue with a list of contributors that are eligible for a reward.
// We do not transmit the name of the issue since it could contain sensible information.
type EligibleIssue struct {
	ID           int64                     `json:"id"`
	Number       int                       `json:"number"`
	HTMLURL      string                    `json:"htmlurl"`
	Contributors []*famedModel.Contributor `json:"contributors"`
}

// GetEligibleIssues returns a list of eligible issues of repos owned by a given owner.
func (gH *githubHandler) GetEligibleIssues(c echo.Context) error {
	var (
		ctx   = c.Request().Context()
		owner = c.Param("owner")
	)

	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, famedModel.ErrMissingOwnerPathParameter.Error())
	}

	response := EligibleIssuesResponse{Repos: []Repo{}}
	repos, err := gH.githubInstallationClient.GetReposByOwner(ctx, owner)
	if err != nil {
		return err
	}

	issueState := model.Closed
	issueOptions := providers.IssueListByRepoOptions{
		Labels: []string{gH.famedConfig.Labels[config.FamedLabelKey].Name},
		State:  &issueState,
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
			if err != nil {
				return err
			}

			if len(contributors) != 0 && err == nil {
				repoResp.EligibleIssues = append(repoResp.EligibleIssues, EligibleIssue{ID: issue.ID, Number: issue.Number, HTMLURL: issue.HTMLURL, Contributors: contributors})
			}
		}

		response.Repos = append(response.Repos, repoResp)
	}

	return c.JSON(http.StatusOK, response)
}
