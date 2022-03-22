package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
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

// GetIssueEvents returns all events for a given issue.
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

type IssuesEvent struct {
	Action string
	Repo
	Issue
}

type Repo struct {
	Name  string
	Owner User
}

func ValidateIssuesEvent(event *github.IssuesEvent) (IssuesEvent, error) {
	var issuesEvent IssuesEvent
	if event.Action == nil ||
		event.Repo == nil ||
		event.Repo.Name == nil ||
		event.Repo.Owner == nil ||
		event.Repo.Owner.Login == nil {
		return issuesEvent, ErrEventMissingData
	}

	switch *event.Action {
	case string(Closed):
		fallthrough

	case string(Assigned):
		fallthrough

	case string(Unassigned):
		fallthrough

	case string(Labeled):
		fallthrough

	case string(Unlabeled):
		issue, err := validateIssue(event.Issue)
		if err != nil {
			return issuesEvent, err
		}

		return IssuesEvent{
			Action: *event.Action,
			Repo: Repo{
				Name: *event.Repo.Name,
				Owner: User{
					Login: *event.Repo.Owner.Login,
				},
			},
			Issue: issue,
		}, nil

	default:
		return issuesEvent, ErrUnhandledEventType
	}
}
