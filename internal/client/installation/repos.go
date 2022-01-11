package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubInstallationClient) GetRepos(ctx context.Context) ([]*github.Repository, error) {
	repoResponse, _, err := c.client.Repositories.List(ctx, c.owner, nil)
	return repoResponse, err
}
