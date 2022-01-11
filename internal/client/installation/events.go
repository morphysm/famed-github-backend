package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func (c *githubInstallationClient) GetRepoEvents(ctx context.Context, repoID string) ([]*github.Event, error) {
	eventsResponse, _, err := c.client.Activity.ListRepositoryEvents(ctx, c.owner, repoID, nil)
	return eventsResponse, err
}
