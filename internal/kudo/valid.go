package kudo

import (
	"log"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

// IsIssueValid checks weather all necessary issue fields are assigned.
func IsIssueValid(issue *github.Issue) bool {
	if issue == nil ||
		issue.ID == nil ||
		issue.Number == nil ||
		issue.Assignee == nil ||
		issue.Assignee.Login == nil ||
		issue.CreatedAt == nil ||
		issue.ClosedAt == nil ||
		issue.Labels == nil {
		log.Printf("[IsIssueValid] missing values in issue with ID: %d", issue.ID)
		return false
	}

	return true
}

// IsValidCloseEvent checks weather all necessary event fields are assigned.
func IsValidCloseEvent(event *github.IssuesEvent, kudoLabel string) bool {
	if !isIssuesEventValid(event) ||
		*event.Action != string(installation.Closed) ||
		!IsIssueValid(event.Issue) ||
		!isIssueKudoLabeled(event.Issue, kudoLabel) {
		log.Println("[IsValidCloseEvent] event is not valid closed event")
		return false
	}

	return true
}

func isIssuesEventValid(event *github.IssuesEvent) bool {
	if event == nil ||
		event.Action == nil ||
		event.Repo == nil ||
		event.Repo.Name == nil {
		log.Println("[isIssuesEventValid] event is not a valid issuesEvent")
		return false
	}

	return true
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
