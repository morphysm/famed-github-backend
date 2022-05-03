package famed

import (
	"context"
	"log"
	"sync"

	"github.com/morphysm/famed-github-backend/internal/config"
	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

type safeWrappedIssue struct {
	sync.RWMutex
	enrichedIssues map[int]model2.EnrichedIssue
}

func (sWI *safeWrappedIssue) Add(wI model2.EnrichedIssue) {
	sWI.Lock()
	defer sWI.Unlock()
	sWI.enrichedIssues[wI.Number] = wI
}

func (gH *githubHandler) loadEnrichedIssues(ctx context.Context, owner string, repoName string) (map[int]model2.EnrichedIssue, error) {
	// Get all issues filtered by label and closed state
	famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
	issueState := model.Closed
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, &issueState)
	if err != nil {
		return nil, err
	}

	return gH.enrichIssues(ctx, owner, repoName, issues), nil
}

func (gH *githubHandler) enrichIssues(ctx context.Context, owner string, repoName string, issues []model.Issue) map[int]model2.EnrichedIssue {
	wg := sync.WaitGroup{}
	safeIssues := safeWrappedIssue{enrichedIssues: make(map[int]model2.EnrichedIssue, len(issues))}
	for _, issue := range issues {
		wg.Add(1)
		go func(ctx context.Context, issue model.Issue) {
			defer wg.Done()

			pullRequest := gH.loadPullRequest(ctx, owner, repoName, issue.Number)
			var events []model.IssueEvent
			if !issue.Migrated {
				events = gH.loadEvents(ctx, owner, repoName, issue.Number)
			}

			enrichedIssue := model2.NewEnrichIssue(issue, pullRequest, events)
			safeIssues.Add(enrichedIssue)
		}(ctx, issue)
	}

	wg.Wait()

	return safeIssues.enrichedIssues
}

func (gH *githubHandler) loadPullRequest(ctx context.Context, owner, repoName string, issueNumber int) *string {
	pullRequest, err := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, issueNumber)
	if pullRequest == nil || err != nil {
		log.Printf("[loadPullRequest] error while requesting pull request for issue with number %d: %v", issueNumber, err)
		return nil
	}
	return pullRequest
}

func (gH *githubHandler) loadEvents(ctx context.Context, owner, repoName string, issueNumber int) []model.IssueEvent {
	events, err := gH.githubInstallationClient.GetIssueEvents(ctx, owner, repoName, issueNumber)
	if err != nil {
		log.Printf("[loadEvents] error while requesting events for issue with number %d: %v", issueNumber, err)
		return nil
	}
	return events
}
