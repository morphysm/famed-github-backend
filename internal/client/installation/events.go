package installation

import (
	"context"
	"log"
	"net/http"
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
	// IssueEventActionHeadRefDeleted & IssueEventActionHeadRefRestored represents a GitHub event action triggered when the pull request’s branch was deleted or restored.
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

// GetIssueEvents returns all events for a given issue.
func (c *githubInstallationClient) GetIssueEvents(ctx context.Context, owner string, repoName string, issueNumber int) ([]IssueEvent, error) {
	var (
		client, _           = c.clients.get(owner)
		allEvents           []*github.IssueEvent
		allCompressedEvents []IssueEvent
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
		compressedEvent, err := validateIssueEvent(event)
		if err != nil {
			continue
		}
		allCompressedEvents = append(allCompressedEvents, compressedEvent)
	}

	return allCompressedEvents, nil
}

func validateIssueEvent(event *github.IssueEvent) (IssueEvent, error) {
	var compressedEvent IssueEvent
	if event == nil ||
		event.ID == nil ||
		event.Event == nil ||
		event.CreatedAt == nil {
		return compressedEvent, ErrEventMissingData
	}

	compressedEvent = IssueEvent{
		ID:        *event.ID,
		Event:     *event.Event,
		CreatedAt: *event.CreatedAt,
		Assignee:  validateUser(event.Assignee),
	}

	return compressedEvent, nil
}

func (c *githubInstallationClient) ValidateWebHookEvent(request *http.Request) (interface{}, error) {
	var event interface{}
	webhookSecret := []byte(c.webhookSecret)

	payload, err := github.ValidatePayload(request, webhookSecret)
	if err != nil {
		return event, err
	}

	event, err = github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return nil, err
	}

	switch event := event.(type) {
	case *github.IssuesEvent:
		issuesEvent, err := validateIssuesEvent(event)
		if err != nil {
			return event, err
		}
		return issuesEvent, err
	case *github.InstallationRepositoriesEvent:
		installationRepositoriesEvent, err := validateInstallationRepositoriesEvent(event)
		if err != nil {
			return event, err
		}
		return installationRepositoriesEvent, err
	case *github.InstallationEvent:
		installationEvent, err := validateInstallationEvent(event)
		if err != nil {
			return event, err
		}
		return installationEvent, err
	default:
		log.Println("[ValidateWebHookEvent] unhandled event")
		return event, ErrUnhandledEventType
	}
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

func validateIssuesEvent(event *github.IssuesEvent) (IssuesEvent, error) {
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

type InstallationRepositoriesEvent struct {
	Action            string
	Installation      RepositoriesInstallation
	RepositoriesAdded []Repository
}

type RepositoriesInstallation struct {
	Account User
}

func validateInstallationRepositoriesEvent(event *github.InstallationRepositoriesEvent) (InstallationRepositoriesEvent, error) {
	var compressedEvent InstallationRepositoriesEvent
	if event == nil || event.Action == nil {
		return compressedEvent, ErrEventMissingData
	}
	if event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil {
		return compressedEvent, ErrEventMissingData
	}
	compressedEvent = InstallationRepositoriesEvent{
		Action: *event.Action,
		Installation: RepositoriesInstallation{
			Account: User{
				Login: *event.Installation.Account.Login,
			},
		},
	}

	for _, repository := range event.RepositoriesAdded {
		compressedEvent.RepositoriesAdded = append(compressedEvent.RepositoriesAdded, Repository{Name: *repository.Name})
	}

	return compressedEvent, nil
}

type InstallationEvent struct {
	Action       string
	Repositories []Repository
	Installation
}

type Installation struct {
	ID      int64
	Account User
}

func validateInstallationEvent(event *github.InstallationEvent) (InstallationEvent, error) {
	var compressedEvent InstallationEvent
	if event == nil ||
		event.Action == nil ||
		event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil ||
		event.Installation.ID == nil {
		return compressedEvent, ErrEventMissingData
	}

	compressedEvent = InstallationEvent{
		Action: *event.Action,
		Installation: Installation{
			ID: *event.Installation.ID,
			Account: User{
				Login: *event.Installation.Account.Login,
			},
		},
	}

	for _, repository := range event.Repositories {
		compressedEvent.Repositories = append(compressedEvent.Repositories, Repository{Name: *repository.Name})
	}

	return compressedEvent, nil
}
