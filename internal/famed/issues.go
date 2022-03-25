package famed

import (
	"context"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

type WrappedIssue struct {
	Issue  github.Issue
	Events []github.IssueEvent
}

func (gH *githubHandler) loadIssuesAndEvents(ctx context.Context, owner string, repoName string) (map[int]WrappedIssue, error) {
	// Get all issues filtered by label and closed state
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, pointer.IssueState(github.Closed))
	if err != nil {
		return nil, echo.ErrBadGateway.SetInternal(err)
	}

	wg := sync.WaitGroup{}
	issues := make(map[int]WrappedIssue, len(issuesResponse))
	for _, issue := range issuesResponse {
		wg.Add(1)

		go func(ctx context.Context, issue github.Issue) {
			defer wg.Done()

			wrappedIssue, err := gH.loadIssueEvents(ctx, owner, repoName, issue)
			if err != nil {
				log.Printf("[loadIssuesAndEvents] error while requesting events for issue with number %d: %v", issue.Number, err)
			}

			issues[issue.Number] = wrappedIssue
		}(ctx, issue)
	}

	wg.Wait()
	return issues, nil
}

func (gH *githubHandler) loadIssueEvents(ctx context.Context, owner string, repoName string, issue github.Issue) (WrappedIssue, error) {
	events, err := gH.githubInstallationClient.GetIssueEvents(ctx, owner, repoName, issue.Number)
	if err != nil {
		return WrappedIssue{}, err
	}

	return WrappedIssue{
		Issue:  issue,
		Events: events,
	}, nil
}
