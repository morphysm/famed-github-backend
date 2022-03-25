package pointer

import "github.com/morphysm/famed-github-backend/internal/client/github"

func String(s string) *string {
	return &s
}

func IssueState(s github.IssueState) *github.IssueState {
	return &s
}
