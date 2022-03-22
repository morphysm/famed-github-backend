package famed

import (
	"errors"
	"log"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var (
	ErrIssueMissingAssignee   = errors.New("the issue is missing an assignee")
	ErrIssueClosedAt          = errors.New("the issue is missing the closed at timestamp")
	ErrIssueMissingData       = errors.New("the issue is missing data promised by the GitHub API")
	ErrIssueMissingFamedLabel = errors.New("the issue is missing the famed label")

	ErrEventMissingData            = errors.New("the event is missing data promised by the GitHub API")
	ErrEventAssigneeMissingData    = errors.New("the event assignee is missing data promised by the GitHub API")
	ErrEventNotClose               = errors.New("the event is not a close event")
	ErrEventNotRepoAdded           = errors.New("the event is not a repo added to installation event")
	ErrEventNotInstallationCreated = errors.New("the event is not a installation created event")
	ErrEventMissingFamedLabel      = errors.New("the event is missing the famed label")
)

// isInstallationEventValid checks weather all necessary event fields are assigned.
func isInstallationEventValid(event *github.InstallationEvent) (bool, error) {
	if event == nil || event.Action == nil {
		return false, ErrEventMissingData
	}
	if *event.Action != "created" {
		return false, ErrEventNotInstallationCreated
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
		return false, ErrEventNotRepoAdded
	}
	if event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil {
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
func isIssueFamedLabeled(issue installation.Issue, famedLabel string) bool {
	for _, label := range issue.Labels {
		if label.Name == famedLabel {
			return true
		}
	}

	log.Printf("[IsIssueFamedLabeled] missing famed label: %s in issue with ID: %d", famedLabel, issue.ID)
	return false
}

func isCommentValid(comment *github.IssueComment) bool {
	return comment != nil && comment.Body != nil
}

func isUserValid(user *github.User) bool {
	return user != nil && user.Login != nil
}
