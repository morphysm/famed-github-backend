package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

type IssueState string

type IssueEventAction string

const (
	Open   IssueState = "open"
	Closed IssueState = "closed"
	All    IssueState = "all"

	// The Actor closed the issue.
	// If the issue was closed by commit message, CommitID holds the SHA1 hash of the commit.
	IssueEventActionClosed IssueEventAction = "closed"
	// The Actor merged into master a branch containing a commit mentioning the issue.
	// CommitID holds the SHA1 of the merge commit.
	IssueEventActionMerged IssueEventAction = "merged"
	// The Actor committed to master a commit mentioning the issue in its commit message.
	// CommitID holds the SHA1 of the commit.
	IssueEventActionReferenced IssueEventAction = "referenced"
	// The Actor did that to the issue.
	IssueEventActionReopened IssueEventAction = "reopened"
	IssueEventActionUnlocked IssueEventAction = "unlocked"
	// The Actor locked the issue.
	// LockReason holds the reason of locking the issue (if provided while locking).
	IssueEventActionLocked IssueEventAction = "locked"
	// The Actor changed the issue title from Rename.From to Rename.To.
	IssueEventActionRenamed IssueEventAction = "renamed"
	// Someone unspecified @mentioned the Actor [sic] in an issue comment body.
	IssueEventActionMentioned IssueEventAction = "mentioned"
	// The Assigner assigned the issue to or removed the assignment from the Assignee.
	IssueEventActionAssigned   IssueEventAction = "assigned"
	IssueEventActionUnassigned IssueEventAction = "unassigned"
	// The Actor added or removed the Label from the issue.
	IssueEventActionLabeled   IssueEventAction = "labeled"
	IssueEventActionUnlabeled IssueEventAction = "unlabeled"
	// The Actor added or removed the issue from the Milestone.
	IssueEventActionMilestoned   IssueEventAction = "milestoned"
	IssueEventActionDemilestoned IssueEventAction = "demilestoned"
	// The Actor subscribed to or unsubscribed from notifications for an issue.
	IssueEventActionSubscribed   IssueEventAction = "subscribed"
	IssueEventActionUnsubscribed IssueEventAction = "unsubscribed"
	// The pull request’s branch was deleted or restored.
	IssueEventActionHeadRefDeleted  IssueEventAction = "head_ref_deleted"
	IssueEventActionHeadRefRestored IssueEventAction = "head_ref_restored"
	// The review was dismissed and `DismissedReview` will be populated below.
	IssueEventActionReviewDismissed IssueEventAction = "review_dismissed"
	// The Actor requested or removed the request for a review.
	// RequestedReviewer and ReviewRequester will be populated below.
	IssueEventActionReviewRequested IssueEventAction = "review_requested"
)

func (c *githubInstallationClient) GetIssuesByRepo(ctx context.Context, repoName string, labels []string, state IssueState) ([]*github.Issue, error) {
	issuesResponse, _, err := c.client.Issues.ListByRepo(ctx, c.owner, repoName, &github.IssueListByRepoOptions{State: string(state), Labels: labels})
	return issuesResponse, err
}

func (c *githubInstallationClient) GetIssueEvents(ctx context.Context, repoName string, issueNumber int) ([]*github.IssueEvent, error) {
	issueEventsResponse, _, err := c.client.Issues.ListIssueEvents(ctx, c.owner, repoName, issueNumber, nil)
	return issueEventsResponse, err
}

func (c *githubInstallationClient) PostComment(ctx context.Context, repoName string, issueNumber int, comment string) (*github.IssueComment, error) {
	issueCommentResponse, _, err := c.client.Issues.CreateComment(ctx, c.owner, repoName, issueNumber, &github.IssueComment{Body: &comment})
	return issueCommentResponse, err
}
