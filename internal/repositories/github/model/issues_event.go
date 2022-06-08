package model

import (
	"github.com/phuslu/log"

	"github.com/google/go-github/v41/github"
)

type IssuesEvent struct {
	Action string
	Repo   Repository
	Issue
}

// NewIssuesEvent validates issue events received through the webhook.
// The changes implied by the event in the case of (un)assigned are mapped to the issue to provide an after event issue state,
// this is done because of a suspected bug in the GitHub API.
func NewIssuesEvent(event *github.IssuesEvent, famedLabel string) (IssuesEvent, error) {
	if event.Action == nil ||
		event.Issue == nil ||
		event.Repo == nil ||
		event.Repo.Name == nil ||
		event.Repo.Owner == nil ||
		event.Repo.Owner.Login == nil {
		return IssuesEvent{}, ErrEventMissingData
	}

	if !isIssueFamedLabeled(event.Issue, famedLabel) {
		return IssuesEvent{}, ErrEventNotFamedLabeled
	}

	switch *event.Action {
	case string(Assigned):
		fallthrough

	case string(Unassigned):
		fallthrough

	case string(Closed):
		fallthrough

	case string(Labeled):
		fallthrough

	case string(Unlabeled):
		// TODO check if this is necessary
		issue, err := NewIssue(event.Issue, *event.Repo.Owner.Login, *event.Repo.Name)
		if err != nil {
			return IssuesEvent{}, err
		}

		owner, err := NewUser(event.Repo.Owner)
		if err != nil {
			return IssuesEvent{}, err
		}

		return IssuesEvent{
			Action: *event.Action,
			Repo: Repository{
				Name:  *event.Repo.Name,
				Owner: owner,
			},
			Issue: issue,
		}, nil

	default:
		return IssuesEvent{}, ErrUnhandledEventType
	}
}

// isIssueFamedLabeled checks weather the issue labels contain expected famed label.
func isIssueFamedLabeled(issue *github.Issue, famedLabel string) bool {
	for _, label := range issue.Labels {
		if *label.Name == famedLabel {
			return true
		}
	}

	log.Warn().Msgf("[IsIssueFamedLabeled] missing famed label: %s in issue with ID: %d", famedLabel, issue.ID)
	return false
}
