package app

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubAppClient) GetInstallations(ctx context.Context) ([]*github.Installation, error) {
	var (
		allInstallations []*github.Installation
		listOptions      = &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}
	)

	for {
		installations, resp, err := c.client.Apps.ListInstallations(ctx, nil)
		if err != nil {
			return allInstallations, err
		}
		allInstallations = append(allInstallations, installations...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	return allInstallations, nil
}
