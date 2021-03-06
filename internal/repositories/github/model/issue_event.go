package model

import (
	"time"

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
	// IssueEventActionHeadRefDeleted & IssueEventActionHeadRefRestored represents a GitHub event action triggered when the pull request???s branch was deleted or restored.
	IssueEventActionHeadRefDeleted  IssueEventAction = "head_ref_deleted"
	IssueEventActionHeadRefRestored IssueEventAction = "head_ref_restored"
	// IssueEventActionReviewDismissed represents a GitHub event action triggered when the review was dismissed and `DismissedReview` will be populated below.
	IssueEventActionReviewDismissed IssueEventAction = "review_dismissed"
	// IssueEventActionReviewRequested represents a GitHub event action triggered when the Actor requested or removed the request for a review.
	// RequestedReviewer and ReviewRequester will be populated below.
	IssueEventActionReviewRequested IssueEventAction = "review_requested"
)

type IssueEvent struct {
	ID        int64
	Event     string
	Assignee  *User
	CreatedAt time.Time
}

func NewIssueEvent(event *github.IssueEvent) (IssueEvent, error) {
	if event == nil ||
		event.ID == nil ||
		event.Event == nil ||
		event.CreatedAt == nil {
		return IssueEvent{}, ErrEventMissingData
	}

	compressedEvent := IssueEvent{
		ID:        *event.ID,
		Event:     *event.Event,
		CreatedAt: *event.CreatedAt,
	}

	assignee, err := NewUser(event.Assignee)
	if err == nil {
		compressedEvent.Assignee = &assignee
	}

	return compressedEvent, nil
}
