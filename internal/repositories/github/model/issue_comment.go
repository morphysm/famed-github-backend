package model

import "github.com/google/go-github/v41/github"

// IssueComment represents a GitHub comment.
type IssueComment struct {
	ID int64
	User
	Body string
}

// NewComment returns a new IssueComment.
// If expected data is not available an error is returned.
func NewComment(comment *github.IssueComment) (IssueComment, error) {
	if comment == nil ||
		comment.Body == nil {
		return IssueComment{}, ErrIssueCommentMissingData
	}

	user, err := NewUser(comment.User)
	if err != nil {
		return IssueComment{}, err
	}

	return IssueComment{
		ID:   *comment.ID,
		User: user,
		Body: *comment.Body,
	}, nil
}
