package kudo

import (
	"errors"
	"log"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

var (
	ErrIssueMissingAssignee = errors.New("the issue is missing an assignee")
	ErrIssueMissingData     = errors.New("the issue is missing data promised by the GitHub API")

	ErrEventMissingData      = errors.New("the event is missing data promised by the GitHub API")
	ErrEventMissingKudoLabel = errors.New("the event is missing the kudo label")
	ErrEventIsNotClose       = errors.New("the event is not a close event")
)

// IsIssueValid checks weather all necessary issue fields are assigned.
func IsIssueValid(issue *github.Issue) (bool, error) {
	if issue.Assignee == nil {
		log.Printf("[IsIssueValid] missing assignee in issue with ID: %d", issue.ID)
		return false, ErrIssueMissingAssignee
	}
	if issue == nil ||
		issue.ID == nil ||
		issue.Number == nil ||
		issue.Assignee.Login == nil ||
		issue.CreatedAt == nil ||
		issue.ClosedAt == nil ||
		issue.Labels == nil {
		log.Printf("[IsIssueValid] missing values in issue with ID: %d", issue.ID)
		return false, ErrIssueMissingData
	}

	return true, nil
}

// IsValidCloseEvent checks weather all necessary event fields are assigned.
func IsValidCloseEvent(event *github.IssuesEvent, kudoLabel string) (bool, error) {
	if _, err := isIssuesEventDataValid(event); err != nil {
		return false, err
	}
	if *event.Action != string(installation.Closed) {
		log.Println("[IsValidCloseEvent] event is not a closed event")
		return false, ErrEventIsNotClose
	}
	if !isIssueKudoLabeled(event.Issue, kudoLabel) {
		log.Println("[IsValidCloseEvent] event is missing the kudo label")
		return false, ErrEventMissingKudoLabel
	}
	if _, err := IsIssueValid(event.Issue); err != nil {
		log.Println("[IsValidCloseEvent] event is missing data")
		return false, err
	}

	return true, nil
}

func isIssuesEventDataValid(event *github.IssuesEvent) (bool, error) {
	if event == nil ||
		event.Action == nil ||
		event.Repo == nil ||
		event.Repo.Name == nil {
		log.Println("[isIssuesEventValid] event is not a valid issuesEvent")
		return false, nil
	}

	return true, nil
}

// isIssueKudoLabeled checks weather the issue labels contain expected kudo label.
func isIssueKudoLabeled(issue *github.Issue, kudoLabel string) bool {
	for _, label := range issue.Labels {
		if label.Name != nil && *label.Name == kudoLabel {
			return true
		}
	}

	log.Printf("[IsIssueKudoLabeled] missing kudo label: %s in issue with ID: %d", kudoLabel, *issue.ID)
	return false
}

func isLabelValid(label *github.Label) bool {
	if label.Name == nil {
		log.Printf("[isLabelValid] missing label name in label with ID: %d", label.ID)
		return false
	}

	return true
}
