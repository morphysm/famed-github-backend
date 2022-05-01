package providers

import (
	"context"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/respositories/github/model"
)

func (c *githubAppClient) GetInstallations(ctx context.Context) ([]model.Installation, error) {
	var (
		allInstallations           []*github.Installation
		allCompressedInstallations []model.Installation
		listOptions                = &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}
	)

	for {
		installations, resp, err := c.client.Apps.ListInstallations(ctx, nil)
		if err != nil {
			return allCompressedInstallations, err
		}
		allInstallations = append(allInstallations, installations...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	for _, installation := range allInstallations {
		compressedInstallation, err := model.NewInstallation(installation)
		if err != nil {
			continue
		}

		allCompressedInstallations = append(allCompressedInstallations, compressedInstallation)
	}

	return allCompressedInstallations, nil
}
