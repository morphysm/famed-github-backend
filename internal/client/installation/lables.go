package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubInstallationClient) GetLabels(ctx context.Context, repoID string) ([]*github.Label, error) {
	labelsResponse, _, err := c.client.Issues.ListLabels(ctx, c.owner, repoID, nil)
	return labelsResponse, err
}
