package famed

import (
	"context"
	"log"
	"sync"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type enrichedIssue struct {
	github.Issue
	PullRequest *string
	Events      []github.IssueEvent
}

func newEnrichIssue(issue github.Issue) enrichedIssue {
	return enrichedIssue{Issue: issue}
}

func (i *enrichedIssue) loadPullRequest(ctx context.Context, gH *githubHandler, owner, repoName string) {
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, i.Number)
	if pullRequest == nil || err != nil {
		log.Printf("[loadPullRequest] error while requesting pull request for issue with number %d: %v", i.Number, err)
		return
	}
	i.PullRequest = &pullRequest.URL
}

func (i *enrichedIssue) loadEvents(ctx context.Context, gH *githubHandler, owner, repoName string) {
	events, err := gH.githubInstallationClient.GetIssueEvents(ctx, owner, repoName, i.Number)
	if err != nil {
		log.Printf("[loadEvents] error while requesting events for issue with number %d: %v", i.Number, err)
		return
	}
	i.Events = events
}

type safeWrappedIssue struct {
	sync.RWMutex
	enrichedIssue map[int]enrichedIssue
}

func (sWI *safeWrappedIssue) Add(wI enrichedIssue) {
	sWI.Lock()
	defer sWI.Unlock()
	sWI.enrichedIssue[wI.Number] = wI
}

func (gH *githubHandler) loadIssues(ctx context.Context, owner string, repoName string) (map[int]enrichedIssue, error) {
	// Get all issues filtered by label and closed state
	famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
	issueState := github.Closed
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, &issueState)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	safeIssues := safeWrappedIssue{enrichedIssue: make(map[int]enrichedIssue, len(issuesResponse))}
	for _, issue := range issuesResponse {
		wg.Add(1)
		go func(ctx context.Context, issue github.Issue) {
			defer wg.Done()

			enrichedIssue := newEnrichIssue(issue)
			enrichedIssue.loadPullRequest(ctx, gH, owner, repoName)
			if !issue.Migrated {
				enrichedIssue.loadEvents(ctx, gH, owner, repoName)
			}

			safeIssues.Add(enrichedIssue)
		}(ctx, issue)
	}

	wg.Wait()
	return safeIssues.enrichedIssue, nil
}
