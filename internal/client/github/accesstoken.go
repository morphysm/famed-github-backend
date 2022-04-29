package github

import (
	"context"

	"github.com/google/go-github/v41/github"
)

// GetAccessToken return an access token for a given installationID.
func (c *githubAppClient) GetAccessToken(ctx context.Context, installationID int64) (*github.InstallationToken, error) {
	token, _, err := c.client.Apps.CreateInstallationToken(
		ctx,
		installationID,
		nil)
	if err != nil {
		return nil, err
	}

	return token, err
}
