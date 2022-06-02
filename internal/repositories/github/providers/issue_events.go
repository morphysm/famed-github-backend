package providers

import (
	"context"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// GetIssueEvents returns all events for a given issue.
func (c *githubInstallationClient) GetIssueEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]model.IssueEvent, error) {
	var (
		client, _           = c.clients.get(owner)
		allEvents           []*github.IssueEvent
		allCompressedEvents []model.IssueEvent
		listOptions         = &github.ListOptions{
			Page:    1,
			PerPage: 100,
		}
	)

	for {
		events, resp, err := client.Issues.ListIssueEvents(ctx, owner, repoName, issueNumber, listOptions)
		if err != nil {
			return allCompressedEvents, err
		}
		allEvents = append(allEvents, events...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, event := range allEvents {
		compressedEvent, err := model.NewIssueEvent(event)
		if err != nil {
			continue
		}
		allCompressedEvents = append(allCompressedEvents, compressedEvent)
	}

	return allCompressedEvents, nil
}
