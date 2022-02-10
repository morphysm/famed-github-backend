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

func (c *githubInstallationClient) GetComments(ctx context.Context, repoName string, issueNumber int) ([]*github.IssueComment, error) {
	// GitHub does not allow get comments in an order (https://docs.github.com/en/rest/reference/issues#list-issue-comments)
	issueCommentResponse, _, err := c.client.Issues.ListComments(ctx, c.owner, repoName, issueNumber, nil)
	return issueCommentResponse, err
}
