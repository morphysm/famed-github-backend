package installation

import (
	"context"

	"github.com/google/go-github/v41/github"
)

type Label struct {
	Name        string
	Color       string
	Description string
}

func (c *githubInstallationClient) PostLabel(ctx context.Context, owner string, repo string, label Label) error {
	client := c.clients[owner]

	_, _, err := client.Issues.CreateLabel(ctx, owner, repo, &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	})
	return err
}
