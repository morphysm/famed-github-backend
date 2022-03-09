package famed

import (
	"errors"
	"log"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var (
	ErrIssueMissingAssignee   = errors.New("the issue is missing an assignee")
	ErrIssueMissingData       = errors.New("the issue is missing data promised by the GitHub API")
	ErrIssueMissingFamedLabel = errors.New("the issue is missing the famed label")

	ErrEventMissingData              = errors.New("the event is missing data promised by the GitHub API")
	ErrEventAssigneeMissingData      = errors.New("the event assignee is missing data promised by the GitHub API")
	ErrEventIsNotClose               = errors.New("the event is not a close event")
	ErrEventIsNotRepoAdded           = errors.New("the event is not a repo added to installation event")
	ErrEventIsNotInstallationCreated = errors.New("the event is not a installation created event")
)

// isIssueValid checks weather all necessary issue fields are assigned.
func isIssueValid(issue *github.Issue) (bool, error) {
	if issue == nil ||
		issue.ID == nil ||
		issue.Number == nil ||
		issue.CreatedAt == nil ||
		issue.ClosedAt == nil {
		log.Printf("[isIssueValid] missing values in issue with ID: %d", issue.ID)
		return false, ErrIssueMissingData
	}
	if issue.Assignee == nil ||
		issue.Assignee.Login == nil {
		log.Printf("[isIssueValid] missing assignee in issue with ID: %d", issue.ID)
		return false, ErrIssueMissingAssignee
	}

	return true, nil
}

func hasIssueNumber(issue *github.Issue) bool {
	return issue != nil && issue.Number != nil
}

func hasEventAction(event *github.IssuesEvent) bool {
	return event.Action != nil
}

func hasEvent(event *github.IssuesEvent) bool {
	return event != nil
}

// isWebhookEventValid checks if the base webhook event fields are assigned.
func isWebhookEventValid(event *github.IssuesEvent) bool {
	return hasEvent(event) && hasEventAction(event) && isRepoValid(event.Repo) && isUserValid(event.Repo.Owner) && hasIssueNumber(event.Issue)
}

// isInstallationEventValid checks weather all necessary event fields are assigned.
func isInstallationEventValid(event *github.InstallationEvent) (bool, error) {
	if event == nil || event.Action == nil {
		return false, ErrEventMissingData
	}
	if *event.Action != "created" {
		return false, ErrEventIsNotInstallationCreated
	}
	if event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil ||
		event.Installation.ID == nil {
		return false, ErrEventMissingData
	}

	return true, nil
}

// isIssueValid checks weather all necessary event fields are assigned.
func isRepoAddedEventValid(event *github.InstallationRepositoriesEvent) (bool, error) {
	if event == nil || event.Action == nil {
		return false, ErrEventMissingData
	}
	if *event.Action != "added" {
		return false, ErrEventIsNotRepoAdded
	}
	if event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil {
		return false, ErrEventMissingData
	}

	return true, nil
}

// isCloseEventValid checks weather all necessary event fields are assigned.
func isCloseEventValid(event *github.IssuesEvent, famedLabel string) (bool, error) {
	if _, err := isIssuesEventDataValid(event); err != nil {
		return false, err
	}
	if *event.Action != string(installation.Closed) {
		log.Println("[isCloseEventValid] event is not a closed event")
		return false, ErrEventIsNotClose
	}
	if !isIssueFamedLabeled(event.Issue, famedLabel) {
		return false, ErrIssueMissingFamedLabel
	}
	if _, err := isIssueValid(event.Issue); err != nil {
		log.Println("[isCloseEventValid] event issue is missing data")
		return false, err
	}

	return true, nil
}

func isIssuesEventDataValid(event *github.IssuesEvent) (bool, error) {
	if event == nil ||
		event.Action == nil ||
		event.Repo == nil ||
		event.Repo.Name == nil ||
		event.Repo.Owner == nil ||
		event.Repo.Owner.Login == nil {
		log.Println("[isIssuesEventValid] event is not a valid issuesEvent")
		return false, ErrEventMissingData
	}

	return true, nil
}

// isIssueUnAssignedEventDataValid checks weather the assigner or unassigned event has all necessary data
func isIssueUnAssignedEventDataValid(event *github.IssueEvent) (bool, error) {
	if event == nil || event.CreatedAt == nil {
		return false, ErrEventMissingData
	}

	return isAssigneeDataValid(event.Assignee)
}

func isAssigneeDataValid(assignee *github.User) (bool, error) {
	if assignee == nil || assignee.Login == nil {
		return false, ErrEventAssigneeMissingData
	}

	return true, nil
}

// isIssueFamedLabeled checks weather the issue labels contain expected famed label.
func isIssueFamedLabeled(issue *github.Issue, famedLabel string) bool {
	for _, label := range issue.Labels {
		if isLabelValid(label) && *label.Name == famedLabel {
			return true
		}
	}

	log.Printf("[IsIssueFamedLabeled] missing famed label: %s in issue with ID: %d", famedLabel, *issue.ID)
	return false
}

func isLabelValid(label *github.Label) bool {
	if label != nil && label.Name == nil {
		log.Printf("[isLabelValid] missing label name in label with ID: %d", label.ID)
		return false
	}

	return true
}

func isCommentValid(comment *github.IssueComment) bool {
	return comment != nil && comment.Body != nil
}

func isUserValid(user *github.User) bool {
	return user != nil && user.Login != nil
}

func isRepoValid(repo *github.Repository) bool {
	return repo != nil && repo.Name != nil
}
