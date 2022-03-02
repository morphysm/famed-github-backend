package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

type IssueState string

const (
	Open   IssueState = "open"
	Closed IssueState = "closed"
	All    IssueState = "all"
)

func (c *githubInstallationClient) GetIssuesByRepo(ctx context.Context, owner string, repoName string, labels []string, state IssueState) ([]*github.Issue, error) {
	var (
		client, _   = c.clients.get(owner)
		allIssues   []*github.Issue
		listOptions = &github.ListOptions{
			Page:    1,
			PerPage: 100,
		}
	)

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repoName, &github.IssueListByRepoOptions{State: string(state), Labels: labels})
		if err != nil {
			return allIssues, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	return allIssues, nil
}
