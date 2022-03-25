package github

import (
	"context"

	"github.com/google/go-github/v41/github"
)

type Installation struct {
	ID      int64
	Account User
}

func (c *githubAppClient) GetInstallations(ctx context.Context) ([]Installation, error) {
	var (
		allInstallations           []*github.Installation
		allCompressedInstallations []Installation
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
		compressedInstallation, err := validateInstallation(installation)
		if err != nil {
			continue
		}

		allCompressedInstallations = append(allCompressedInstallations, compressedInstallation)
	}

	return allCompressedInstallations, nil
}

func validateInstallation(installation *github.Installation) (Installation, error) {
	if installation == nil ||
		installation.ID == nil {
		return Installation{}, ErrInstallationMissingData
	}

	account, err := validateUser(installation.Account)
	if err != nil {
		return Installation{}, err
	}

	return Installation{ID: *installation.ID, Account: account}, nil
}
