package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/pointers"
)

type IssueState string

const (
	Opened     IssueState = "opened"
	Closed     IssueState = "closed"
	Reopened   IssueState = "reopened"
	Edited     IssueState = "edited"
	Assigned   IssueState = "assigned"
	Unassigned IssueState = "unassigned"
	Labeled    IssueState = "labeled"
	Unlabeled  IssueState = "unlabeled"
)

func (c *githubInstallationClient) GetIssuesByRepo(ctx context.Context, owner string, repoName string, labels []string, state IssueState) ([]*github.Issue, error) {
	var (
		client, _   = c.clients.get(owner)
		allIssues   []*github.Issue
		listOptions = &github.ListOptions{
			Page:    1,
			PerPage: 100,
		}
	)

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repoName, &github.IssueListByRepoOptions{State: string(state), Labels: labels})
		if err != nil {
			return allIssues, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	return allIssues, nil
}

// ReopenIssue reopens a given issue.
func (c *githubInstallationClient) ReopenIssue(ctx context.Context, owner string, repoName string, issueNumber int) error {
	var (
		client, _ = c.clients.get(owner)
		issue     = &github.IssueRequest{
			State: pointers.String(string(Reopened)),
		}
	)

	_, _, err := client.Issues.Edit(ctx, owner, repoName, issueNumber, issue)
	if err != nil {
		return err
	}

	return nil
}

// GetIssuePullRequest returns a pull request if a linked pull request for the given issue can be found.
// This is a workaround for the missing "pull_request" field in the event and issue objects provided by the REST GitHub API.
// https://github.community/t/get-referenced-pull-request-from-issue/14027
func (c *githubInstallationClient) GetIssuePullRequest(ctx context.Context, owner string, repoName string, issueNumber int) (*PullRequest, error) {
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

	return &lastConnectedEvent.ConnectedEvent.Subject.PullRequest, nil
}
