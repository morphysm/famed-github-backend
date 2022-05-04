package providers

import (
	"context"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
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

// DeleteComment delets a given GitHub comment.
func (c *githubInstallationClient) DeleteComment(ctx context.Context, owner string, repoName string, commentID int64) error {
	client, _ := c.clients.get(owner)

	_, err := client.Issues.DeleteComment(ctx, owner, repoName, commentID)
	return err
}

// GetComments returns all GitHub comments of a given GitHub issue.
func (c *githubInstallationClient) GetComments(ctx context.Context, owner string, repoName string, issueNumber int) ([]model.IssueComment, error) {
	// GitHub does not allow get comments in an order (https://docs.github.com/en/rest/reference/issues#list-issue-comments)
	var (
		client, _             = c.clients.get(owner)
		allComments           []*github.IssueComment
		allCompressedComments []model.IssueComment
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
		compressedComment, err := model.NewComment(comment)
		if err != nil {
			continue
		}
		allCompressedComments = append(allCompressedComments, compressedComment)
	}

	return allCompressedComments, nil
}
