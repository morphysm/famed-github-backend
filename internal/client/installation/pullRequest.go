package installation

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
		PullRequest PullRequest `graphql:"... on PullRequest"`
	}
	CreatedAt time.Time
}

type PullRequest struct {
	URL string
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
