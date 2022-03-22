package famed

import (
	"errors"
	"log"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var (
	ErrIssueMissingAssignee = errors.New("the issue is missing an assignee")
	ErrIssueClosedAt        = errors.New("the issue is missing the closed at timestamp")

	ErrEventMissingData            = errors.New("the event is missing data promised by the GitHub API")
	ErrEventAssigneeMissingData    = errors.New("the event assignee is missing data promised by the GitHub API")
	ErrEventNotRepoAdded           = errors.New("the event is not a repo added to installation event")
	ErrEventNotInstallationCreated = errors.New("the event is not a installation created event")
	ErrEventMissingFamedLabel      = errors.New("the event is missing the famed label")
)

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
