package github

import (
	"context"
	"errors"
	"net/http"
	"sort"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

func (gH *githubHandler) GetContributors(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo name path parameter"))
	}

	// Get all issues in repo
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), repoName, []string{gH.kudoLabel}, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	// Use issues to generate contributor list
	contributors, err := gH.issuesToContributors(c.Request().Context(), issuesResponse, repoName)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}

// TODO test if issues are returned in chronological order
func (gH *githubHandler) issuesToContributors(ctx context.Context, issues []*github.Issue, repoName string) ([]*kudo.Contributor, error) {
	var (
		contributors      map[string]*kudo.Contributor
		contributorsArray []*kudo.Contributor
	)

	for _, issue := range issues {
		if issue.ID == nil || issue.CreatedAt == nil || issue.ClosedAt == nil {
			continue
		}

		eventsResp, err := gH.githubInstallationClient.GetIssueEvents(ctx, repoName, *issue.Number)
		if err != nil {
			return nil, err
		}

		severity := kudo.IssueToSeverity(issue)

		contributors = kudo.EventsToContributors(contributors, eventsResp, *issue.CreatedAt, *issue.ClosedAt, severity)
	}

	// Transformation of contributors map to contributors array
	for _, contributor := range contributors {
		contributorsArray = append(contributorsArray, contributor)
	}

	// Sort contributors array by total rewards
	sort.SliceStable(contributorsArray, func(i, j int) bool {
		return contributorsArray[i].RewardSum > contributorsArray[j].RewardSum
	})

	return contributorsArray, nil
}
