package github

import (
	"context"

	"github.com/google/go-github/v41/github"
)

// PostComment posts a comment to a given GitHub issue.
func (c *githubInstallationClient) PostComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string) error {
	client, _ := c.clients.get(owner)

	_, _, err := client.Issues.CreateComment(ctx, owner, repoName, issueNumber, &github.IssueComment{Body: &comment})
	return err
}

// UpdateComment updates a given GitHub comment.
func (c *githubInstallationClient) UpdateComment(ctx context.Context, owner string, repoName string, commentID int64, comment string) error {
	client, _ := c.clients.get(owner)

	_, _, err := client.Issues.EditComment(ctx, owner, repoName, commentID, &github.IssueComment{Body: &comment})
	return err
}

// IssueComment represents a GitHub comment.
type IssueComment struct {
	ID int64
	User
	Body string
}

// GetComments returns all GitHub comments of a given GitHub issue.
func (c *githubInstallationClient) GetComments(ctx context.Context, owner string, repoName string, issueNumber int) ([]IssueComment, error) {
	// GitHub does not allow get comments in an order (https://docs.github.com/en/rest/reference/issues#list-issue-comments)
	var (
		client, _             = c.clients.get(owner)
		allComments           []*github.IssueComment
		allCompressedComments []IssueComment
		listOptions           = &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}
	)

	for {
		comments, resp, err := client.Issues.ListComments(ctx, owner, repoName, issueNumber, listOptions)
		if err != nil {
			return allCompressedComments, err
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, comment := range allComments {
		compressedComment, err := validateComment(comment)
		if err != nil {
			continue
		}
		allCompressedComments = append(allCompressedComments, compressedComment)
	}

	return allCompressedComments, nil
}

// GetComments returns a compressed IssueComment.
// If expected data is not available an error is returned.
func validateComment(comment *github.IssueComment) (IssueComment, error) {
	if comment == nil ||
		comment.Body == nil {
		return IssueComment{}, ErrIssueCommentMissingData
	}

	user, err := validateUser(comment.User)
	if err != nil {
		return IssueComment{}, err
	}

	return IssueComment{
		ID:   *comment.ID,
		User: user,
		Body: *comment.Body,
	}, nil
}
