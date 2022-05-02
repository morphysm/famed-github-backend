package model

import (
	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

type EnrichedIssue struct {
	model.Issue
	PullRequest *string
	Events      []model.IssueEvent
}

func NewEnrichIssue(issue model.Issue, pullRequest *string, events []model.IssueEvent) EnrichedIssue {
	return EnrichedIssue{
		Issue:       issue,
		PullRequest: pullRequest,
		Events:      events,
	}
}
