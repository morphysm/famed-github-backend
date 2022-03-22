package installation

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v41/github"
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

type Issue struct {
	ID        int64
	Number    int
	Title     string
	CreatedAt time.Time
	ClosedAt  *time.Time
	Assignee  *User
	Labels    []Label
}

func validateIssue(issue *github.Issue) (Issue, error) {
	var compressedIssue Issue
	if issue == nil ||
		issue.ID == nil ||
		issue.Number == nil ||
		issue.Title == nil ||
		issue.CreatedAt == nil {
		return compressedIssue, ErrIssueMissingData
	}

	compressedIssue = Issue{
		ID:        *issue.ID,
		Number:    *issue.Number,
		Title:     *issue.Title,
		CreatedAt: *issue.CreatedAt,
		ClosedAt:  issue.ClosedAt,
		Assignee:  validateUser(issue.Assignee),
	}

	if issue.Labels != nil {
		for _, label := range issue.Labels {
			if label.Name == nil {
				continue
			}
			compressedLabel := Label{Name: *label.Name}
			compressedIssue.Labels = append(compressedIssue.Labels, compressedLabel)
		}
	}

	return compressedIssue, nil
}

type User struct {
	Login      string
	AvatarURL  *string
	HTMLURL    *string
	GravatarID *string
}

func validateUser(user *github.User) *User {
	if user != nil &&
		user.Login != nil {
		return &User{
			Login:      *user.Login,
			AvatarURL:  user.AvatarURL,
			HTMLURL:    user.HTMLURL,
			GravatarID: user.GravatarID,
		}
	}

	return nil
}

func (c *githubInstallationClient) GetIssuesByRepo(ctx context.Context, owner string, repoName string, labels []string, state IssueState) ([]Issue, error) {
	var (
		client, _           = c.clients.get(owner)
		allIssues           []*github.Issue
		allCompressedIssues []Issue
		listOptions         = &github.ListOptions{
			Page:    1,
			PerPage: 100,
		}
	)

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repoName, &github.IssueListByRepoOptions{State: string(state), Labels: labels})
		if err != nil {
			return allCompressedIssues, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, issue := range allIssues {
		compressedIssue, err := validateIssue(issue)
		if err != nil {
			log.Printf("validation error for issue with number %d: %v", issue.Number, err)
		}

		allCompressedIssues = append(allCompressedIssues, compressedIssue)
	}

	return allCompressedIssues, nil
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
