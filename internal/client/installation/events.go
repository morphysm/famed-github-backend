package installation

import (
	"context"
	"sync"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/shurcooL/githubv4"
)

type IssueEventAction string

const (
	// IssueEventActionClosed represents a GitHub event action triggered when the Actor closed the issue.
	// If the issue was closed by commit message, CommitID holds the SHA1 hash of the commit.
	IssueEventActionClosed IssueEventAction = "closed"
	// IssueEventActionMerged represents a GitHub event action triggered when the Actor merged into master a branch containing a commit mentioning the issue.
	// CommitID holds the SHA1 of the merge commit.
	IssueEventActionMerged IssueEventAction = "merged"
	// IssueEventActionReferenced represents a GitHub event action triggered when the Actor committed to master a commit mentioning the issue in its commit message.
	// CommitID holds the SHA1 of the commit.
	IssueEventActionReferenced IssueEventAction = "referenced"
	// IssueEventActionReopened represents a GitHub event action triggered when the Actor reopened the issue.
	IssueEventActionReopened IssueEventAction = "reopened"
	// IssueEventActionUnlocked represents a GitHub event action triggered when the Actor unlocked the issue.
	IssueEventActionUnlocked IssueEventAction = "unlocked"
	// IssueEventActionLocked represents a GitHub event action triggered when the Actor locked the issue.
	// LockReason holds the reason of locking the issue (if provided while locking).
	IssueEventActionLocked IssueEventAction = "locked"
	// IssueEventActionRenamed represents a GitHub event action triggered when the Actor changed the issue title from Rename.
	// From to Rename.To.
	IssueEventActionRenamed IssueEventAction = "renamed"
	// IssueEventActionMentioned represents a GitHub event action triggered when someone unspecified @mentioned the Actor [sic] in an issue comment body.
	IssueEventActionMentioned IssueEventAction = "mentioned"
	// IssueEventActionAssigned represents a GitHub event action triggered when the Assigner assigned the issue to or removed the assignment from the Assignee.
	IssueEventActionAssigned   IssueEventAction = "assigned"
	IssueEventActionUnassigned IssueEventAction = "unassigned"
	// IssueEventActionLabeled & IssueEventActionUnlabeled represents a GitHub event action triggered when the Actor added or removed the Label from the issue.
	IssueEventActionLabeled   IssueEventAction = "labeled"
	IssueEventActionUnlabeled IssueEventAction = "unlabeled"
	// IssueEventActionMilestoned & IssueEventActionDemilestoned represents a GitHub event action triggered when the Actor added or removed the issue from the Milestone.
	IssueEventActionMilestoned   IssueEventAction = "milestoned"
	IssueEventActionDemilestoned IssueEventAction = "demilestoned"
	// IssueEventActionSubscribed & IssueEventActionUnsubscribed represents a GitHub event action triggered when the Actor subscribed to or unsubscribed from notifications for an issue.
	IssueEventActionSubscribed   IssueEventAction = "subscribed"
	IssueEventActionUnsubscribed IssueEventAction = "unsubscribed"
	// IssueEventActionHeadRefDeleted & IssueEventActionHeadRefRestored represents a GitHub event action triggered when the pull request’s branch was deleted or restored.
	IssueEventActionHeadRefDeleted  IssueEventAction = "head_ref_deleted"
	IssueEventActionHeadRefRestored IssueEventAction = "head_ref_restored"
	// IssueEventActionReviewDismissed represents a GitHub event action triggered when the review was dismissed and `DismissedReview` will be populated below.
	IssueEventActionReviewDismissed IssueEventAction = "review_dismissed"
	// IssueEventActionReviewRequested represents a GitHub event action triggered when the Actor requested or removed the request for a review.
	// RequestedReviewer and ReviewRequester will be populated below.
	IssueEventActionReviewRequested IssueEventAction = "review_requested"
)

func (c *githubInstallationClient) GetIssueEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]*github.IssueEvent, error) {
	var (
		client, _   = c.clients.get(owner)
		allEvents   []*github.IssueEvent
		listOptions = &github.ListOptions{
			Page:    1,
			PerPage: 100,
		}
	)

	for {
		events, resp, err := client.Issues.ListIssueEvents(ctx, owner, repoName, issueNumber, listOptions)
		if err != nil {
			return allEvents, err
		}
		allEvents = append(allEvents, events...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	return allEvents, nil
}

// GetIssuesEvents requests all events of a set of issues from the GitHub API in a concurrent fashion.
func (c *githubInstallationClient) GetIssuesEvents(ctx context.Context, owner string, repoName string, issues []*github.Issue) (map[int][]*github.IssueEvent, map[int]error) {
	var (
		allEvents = make(map[int][]*github.IssueEvent)
		errors    = make(map[int]error)
	)

	// Create context with cancel to cancel all request if one fails
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create wait group to wait for all requests to finish
	wg := sync.WaitGroup{}

	for _, issue := range issues {
		wg.Add(1)

		// Start go routine to get the issue's events
		go func(ctx context.Context, repoName string, issueNumber int, issueID int64) {
			defer wg.Done()

			// Check if one of the requests returned an error otherwise run the request
			select {
			case <-ctx.Done():
				return
			default:
				events, err := c.GetIssueEvents(ctx, owner, repoName, issueNumber)
				if err != nil {
					errors[issueNumber] = err
					return
				}
				allEvents[issueNumber] = events
			}
		}(ctx, repoName, *issue.Number, *issue.ID)
	}

	wg.Wait()
	return allEvents, errors
}

type IssueTimelineDisconnectionItem struct {
	DisconnectedEvent ConnectedEvent `graphql:"... on DisconnectedEvent"`
}

type IssueTimelineConnectionItem struct {
	ConnectedEvent ConnectedEvent `graphql:"... on ConnectedEvent"`
}

type ConnectedEvent struct {
	Subject struct {
		PullRequest PullRequest `graphql:"... on PullRequest"`
	}
	CreatedAt time.Time
}

type PullRequest struct {
	URL string
}

func (c *githubInstallationClient) GetDisconnectedEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]IssueTimelineDisconnectionItem, error) {
	var (
		client, _        = c.clients.getGql(owner)
		allTimelineItems []IssueTimelineDisconnectionItem
		query            struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Nodes    []IssueTimelineDisconnectionItem
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

func (c *githubInstallationClient) GetConnectedEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]IssueTimelineConnectionItem, error) {
	var (
		client, _        = c.clients.getGql(owner)
		allTimelineItems []IssueTimelineConnectionItem
		query            struct {
			Repository struct {
				Issue struct {
					TimelineItems struct {
						Nodes    []IssueTimelineConnectionItem
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
