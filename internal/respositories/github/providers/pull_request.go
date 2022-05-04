package providers

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type issueTimelineDisconnectionItem struct {
	DisconnectedEvent connectedEvent `graphql:"... on DisconnectedEvent"`
}

type issueTimelineConnectionItem struct {
	ConnectedEvent connectedEvent `graphql:"... on ConnectedEvent"`
}

type connectedEvent struct {
	Subject struct {
		PullRequest pullRequest `graphql:"... on PullRequest"`
	}
	CreatedAt time.Time
}

type pullRequest struct {
	URL string
}

// GetIssuePullRequest returns a pull request if a linked pull request for the given issue can be found.
// This is a workaround for the missing "pull_request" field in the event and issue objects provided by the REST GitHub API.
// https://github.community/t/get-referenced-pull-request-from-issue/14027
func (c *githubInstallationClient) GetIssuePullRequest(ctx context.Context, owner string, repoName string, issueNumber int) (*string, error) {
	allTimelineItemsConnected, err := c.getConnectedEvents(ctx, owner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	var lastConnectedEvent *issueTimelineConnectionItem
	for _, node := range allTimelineItemsConnected {
		if node.ConnectedEvent.Subject.PullRequest.URL != "" &&
			(lastConnectedEvent == nil || lastConnectedEvent.ConnectedEvent.CreatedAt.Before(node.ConnectedEvent.CreatedAt)) {
			// Last pull request connected event
			tmpN := node
			lastConnectedEvent = &tmpN
		}
	}

	if lastConnectedEvent == nil {
		return nil, nil
	}

	allTimelineItemsDisconnected, err := c.getDisconnectedEvents(ctx, owner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	for _, node := range allTimelineItemsDisconnected {
		if node.DisconnectedEvent.Subject.PullRequest.URL != "" &&
			lastConnectedEvent.ConnectedEvent.CreatedAt.Before(node.DisconnectedEvent.CreatedAt) {
			// Pull request disconnected after last connected event
			return nil, nil
		}
	}

	return &lastConnectedEvent.ConnectedEvent.Subject.PullRequest.URL, nil
}

// getDisconnectedEvents returns all IssueTimelineDisconnectionItems for a given issue.
// This is used as a workaround for the missing "pull_request" field in the event and issue objects provided by the REST GitHub API.
func (c *githubInstallationClient) getDisconnectedEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]issueTimelineDisconnectionItem, error) {
	var (
		client, _        = c.clients.getGql(owner)
		allTimelineItems []issueTimelineDisconnectionItem
		query            struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Nodes    []issueTimelineDisconnectionItem
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"timelineItems(first: 100, after: $commentsCursor)"`
				} `graphql:"issue(number: $issueNumber)"`
			} `graphql:"repository(owner: $owner, name: $repoName)"`
		}
		variables = map[string]interface{}{
			"owner":          githubv4.String(owner),
			"repoName":       githubv4.String(repoName),
			"issueNumber":    githubv4.Int(issueNumber),
			"commentsCursor": (*githubv4.String)(nil),
		}
	)

	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}

		allTimelineItems = append(allTimelineItems, query.Repository.Issue.TimelineItems.Nodes...)
		if !query.Repository.Issue.TimelineItems.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = githubv4.NewString(query.Repository.Issue.TimelineItems.PageInfo.EndCursor)
	}

	return allTimelineItems, nil
}

// getConnectedEvents returns all issueTimelineConnectionItem for a given issue.
// This is used as a workaround for the missing "pull_request" field in the event and issue objects provided by the REST GitHub API.
func (c *githubInstallationClient) getConnectedEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]issueTimelineConnectionItem, error) {
	var (
		client, _        = c.clients.getGql(owner)
		allTimelineItems []issueTimelineConnectionItem
		query            struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Nodes    []issueTimelineConnectionItem
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"timelineItems(first: 100, after: $commentsCursor)"`
				} `graphql:"issue(number: $issueNumber)"`
			} `graphql:"repository(owner: $owner, name: $repoName)"`
		}
		variables = map[string]interface{}{
			"owner":          githubv4.String(owner),
			"repoName":       githubv4.String(repoName),
			"issueNumber":    githubv4.Int(issueNumber),
			"commentsCursor": (*githubv4.String)(nil),
		}
	)

	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}

		allTimelineItems = append(allTimelineItems, query.Repository.Issue.TimelineItems.Nodes...)
		if !query.Repository.Issue.TimelineItems.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = githubv4.NewString(query.Repository.Issue.TimelineItems.PageInfo.EndCursor)
	}

	return allTimelineItems, nil
}
