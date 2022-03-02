package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubInstallationClient) PostComment(ctx context.Context, owner string, repoName string, issueNumber int, comment string) error {
	client, _ := c.clients.get(owner)

	_, _, err := client.Issues.CreateComment(ctx, owner, repoName, issueNumber, &github.IssueComment{Body: &comment})
	return err
}

func (c *githubInstallationClient) GetComments(ctx context.Context, owner string, repoName string, issueNumber int) ([]*github.IssueComment, error) {
	// GitHub does not allow get comments in an order (https://docs.github.com/en/rest/reference/issues#list-issue-comments)
	var (
		client, _   = c.clients.get(owner)
		allComments []*github.IssueComment
		listOptions = &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}
	)

	for {
		comments, resp, err := client.Issues.ListComments(ctx, owner, repoName, issueNumber, listOptions)
		if err != nil {
			return allComments, err
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	return allComments, nil
}
