package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

type IssueState string

const (
	Opened     IssueState = "opened"
	Closed     IssueState = "closed"
	Reopened   IssueState = "reopened"
	Edited     IssueState = "edited"
	Assigned   IssueState = "assigned"
	Unassigned IssueState = "unassigned"
	Labeled    IssueState = "labeled"
	Unlabeled  IssueState = "unlabeled"
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
