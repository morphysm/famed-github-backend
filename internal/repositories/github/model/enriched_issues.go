package model

type EnrichedIssue struct {
	Issue
	PullRequest *string
	Events      []IssueEvent
}

func NewEnrichIssue(issue Issue, pullRequest *string, events []IssueEvent) EnrichedIssue {
	return EnrichedIssue{
		Issue:       issue,
		PullRequest: pullRequest,
		Events:      events,
	}
}
