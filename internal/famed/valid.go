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

	ErrEventMissingData         = errors.New("the event is missing data promised by the GitHub API")
	ErrEventAssigneeMissingData = errors.New("the event assignee is missing data promised by the GitHub API")
	ErrEventIsNotClose          = errors.New("the event is not a close event")
	ErrEventIsNotRepoAdded      = errors.New("the event is not a repo added to installation event")
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

// isValidInstallationEvent checks weather all necessary event fields are assigned.
func isValidInstallationEvent(event *github.InstallationEvent) (bool, error) {
	if event == nil || event.Action == nil {
		return false, ErrEventMissingData
	}
	if *event.Action != "created" {
		return false, ErrEventIsNotRepoAdded
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
func isValidRepoAddedEvent(event *github.InstallationRepositoriesEvent) (bool, error) {
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

// isValidCloseEvent checks weather all necessary event fields are assigned.
func isValidCloseEvent(event *github.IssuesEvent, famedLabel string) (bool, error) {
	if _, err := isIssuesEventDataValid(event); err != nil {
		return false, err
	}
	if *event.Action != string(installation.Closed) {
		log.Println("[isValidCloseEvent] event is not a closed event")
		return false, ErrEventIsNotClose
	}
	if !isIssueFamedLabeled(event.Issue, famedLabel) {
		return false, ErrIssueMissingFamedLabel
	}
	if _, err := isIssueValid(event.Issue); err != nil {
		log.Println("[isValidCloseEvent] event issue is missing data")
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
		if label.Name != nil && *label.Name == famedLabel {
			return true
		}
	}

	log.Printf("[IsIssueFamedLabeled] missing famed label: %s in issue with ID: %d", famedLabel, *issue.ID)
	return false
}

func isLabelValid(label *github.Label) bool {
	if label.Name == nil {
		log.Printf("[isLabelValid] missing label name in label with ID: %d", label.ID)
		return false
	}

	return true
}

func isCommentValid(comment *github.IssueComment) bool {
	return comment != nil && comment.Body != nil
}
