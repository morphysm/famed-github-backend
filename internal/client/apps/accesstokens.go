package apps

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubAppClient) GetAccessTokens(ctx context.Context, installationID int64, repositoryIDs []int64) (*github.InstallationToken, error) {
	token, _, err := c.client.Apps.CreateInstallationToken(
		ctx,
		installationID,
		&github.InstallationTokenOptions{RepositoryIDs: repositoryIDs})
	if err != nil {
		return nil, err
	}

	return token, err
}
