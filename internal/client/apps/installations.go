package apps

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubAppClient) GetInstallations(ctx context.Context) ([]*github.Installation, error) {
	installationResponse, _, err := c.client.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return installationResponse, err
}
