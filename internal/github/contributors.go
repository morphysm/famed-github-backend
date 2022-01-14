package github

import (
	"context"
	"errors"
	"log"
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
		contributorsArray []*kudo.Contributor
		filteredIssues    []*github.Issue
		eventsByIssue     = map[int64][]*github.IssueEvent{}
	)

	for _, issue := range issues {
		if !kudo.IsIssueValid(issue) {
			log.Printf("[issuesToContributors] issue invalid with ID: %d \n", issue.ID)
			continue
		}
		filteredIssues = append(filteredIssues, issue)

		// TODO add concurrency
		eventsResp, err := gH.githubInstallationClient.GetIssueEvents(ctx, repoName, *issue.Number)
		if err != nil {
			return nil, err
		}

		eventsByIssue[*issue.ID] = eventsResp
	}

	contributors := kudo.GenerateContributors(filteredIssues, eventsByIssue)

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
