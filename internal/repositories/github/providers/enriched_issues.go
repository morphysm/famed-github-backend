package providers

import (
	"context"
	"sync"

	"github.com/phuslu/log"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type safeWrappedIssue struct {
	sync.RWMutex
	enrichedIssues map[int]model.EnrichedIssue
}

func (sWI *safeWrappedIssue) Add(wI model.EnrichedIssue) {
	sWI.Lock()
	defer sWI.Unlock()
	sWI.enrichedIssues[wI.Number] = wI
}

func (c *githubInstallationClient) GetEnrichedIssues(ctx context.Context, owner string, repoName string) (map[int]model.EnrichedIssue, error) {
	issueState := model.Closed
	issues, err := c.GetIssuesByRepo(ctx, owner, repoName, []string{c.famedLabel}, &issueState)
	if err != nil {
		return nil, err
	}

	return c.EnrichIssues(ctx, owner, repoName, issues), nil
}

func (c *githubInstallationClient) EnrichIssues(ctx context.Context, owner string, repoName string, issues []model.Issue) map[int]model.EnrichedIssue {
	wg := sync.WaitGroup{}
	safeIssues := safeWrappedIssue{enrichedIssues: make(map[int]model.EnrichedIssue, len(issues))}
	for _, issue := range issues {
		wg.Add(1)
		go func(ctx context.Context, owner string, repoName string, issue model.Issue) {
			defer wg.Done()

			enrichedIssue := c.EnrichIssue(ctx, owner, repoName, issue)
			safeIssues.Add(enrichedIssue)
		}(ctx, owner, repoName, issue)
	}

	wg.Wait()
	return safeIssues.enrichedIssues
}

func (c *githubInstallationClient) EnrichIssue(ctx context.Context, owner string, repoName string, issue model.Issue) model.EnrichedIssue {
	pullRequest, err := c.GetIssuePullRequest(ctx, owner, repoName, issue.Number)
	if pullRequest == nil || err != nil {
		log.Error().Err(err).Msgf("[EnrichIssue] error while requesting pull request for issue with number %d", issue.Number)
	}

	var events []model.IssueEvent
	if !issue.Migrated {
		events, err = c.GetIssueEvents(ctx, owner, repoName, issue.Number)
		if err != nil {
			log.Error().Err(err).Msgf("[EnrichIssue] error while requesting events for issue with number %d", issue.Number)
		}
	}

	enrichedIssue := model.NewEnrichIssue(issue, pullRequest, events)
	return enrichedIssue
}
