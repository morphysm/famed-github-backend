package famed

import (
	"context"
	"log"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type WrappedIssue struct {
	Issue  *github.Issue
	Events []*github.IssueEvent
}

func (gH *githubHandler) loadIssuesAndEvents(ctx context.Context, owner string, repoName string) (map[int]WrappedIssue, error) {
	// Get all issues filtered by label and closed state
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, installation.Closed)
	if err != nil {
		return nil, echo.ErrBadGateway.SetInternal(err)
	}

	filteredIssues := filterIssues(issuesResponse)
	if len(filteredIssues) == 0 {
		return nil, nil
	}

	wg := sync.WaitGroup{}
	issues := make(map[int]WrappedIssue, len(filteredIssues))
	for _, issue := range filteredIssues {
		wg.Add(1)

		go func(ctx context.Context, issue *github.Issue) {
			defer wg.Done()

			wrappedIssue, err := gH.loadIssueEvents(ctx, owner, repoName, issue)
			if err != nil {
				log.Printf("[loadIssuesAndEvents] error while requesting events for issue with number %d: %v", issue.Number, err)
			}

			issues[*issue.Number] = wrappedIssue
		}(ctx, issue)
	}

	wg.Wait()
	return issues, nil
}

func (gH *githubHandler) loadIssueEvents(ctx context.Context, owner string, repoName string, issue *github.Issue) (WrappedIssue, error) {
	var wrappedIssue WrappedIssue

	events, err := gH.githubInstallationClient.GetIssueEvents(ctx, owner, repoName, *issue.Number)
	if err != nil {
		return wrappedIssue, err
	}

	wrappedIssue.Issue = issue
	wrappedIssue.Events = events

	return wrappedIssue, nil
}

// filterIssues filters for valid issues.
func filterIssues(issues []*github.Issue) []*github.Issue {
	filteredIssues := make([]*github.Issue, 0)
	for _, issue := range issues {
		if _, err := isIssueValid(issue); err != nil {
			log.Printf("[issuesToContributors] issue invalid with ID: %d, error: %v \n", issue.ID, err)
			continue
		}
		filteredIssues = append(filteredIssues, issue)
	}

	return filteredIssues
}
