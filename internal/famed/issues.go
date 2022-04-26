package famed

import (
	"context"
	"log"
	"sync"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/internal/config"
)

type WrappedIssue struct {
	Issue       github.Issue
	PullRequest *github.PullRequest
	Events      []github.IssueEvent
}

type safeWrappedIssue struct {
	sync.RWMutex
	wrappedIssues map[int]WrappedIssue
}

func (sWI *safeWrappedIssue) Add(wI WrappedIssue) {
	sWI.Lock()
	defer sWI.Unlock()
	sWI.wrappedIssues[wI.Issue.Number] = wI
}

func (gH *githubHandler) loadIssues(ctx context.Context, owner string, repoName string) (map[int]WrappedIssue, error) {
	// Get all issues filtered by label and closed state
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]
	issueState := github.Closed
	issuesResponse, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, owner, repoName, []string{famedLabel.Name}, &issueState)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	safeIssues := safeWrappedIssue{wrappedIssues: make(map[int]WrappedIssue, len(issuesResponse))}
	for _, issue := range issuesResponse {
		// TODO Skip event loading for migrated issues because information is already present
		//if issue.Migrated {
		//	issues[issue.Number] = WrappedIssue{Issue: issue}
		//	continue
		//}

		// TODO refactor
		// TODO commented out for DevConnect
		pullRequest, _ := gH.githubInstallationClient.GetIssuePullRequest(ctx, owner, repoName, issue.Number)
		//if pullRequest == nil || err != nil {
		//	safeIssues.Add(WrappedIssue{Issue: issue, PullRequest: nil})
		//	continue
		//}

		wg.Add(1)

		go func(ctx context.Context, issue github.Issue, pullRequest *github.PullRequest) {
			defer wg.Done()

			wrappedIssue, err := gH.loadIssueEvents(ctx, owner, repoName, issue)
			if err != nil {
				log.Printf("[loadIssues] error while requesting events for issue with number %d: %v", issue.Number, err)
			}

			wrappedIssue.PullRequest = pullRequest
			safeIssues.Add(wrappedIssue)
		}(ctx, issue, pullRequest)
	}

	wg.Wait()
	return safeIssues.wrappedIssues, nil
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
