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
	client, _ := c.clients.get(owner)

	_, _, err := client.Issues.CreateLabel(ctx, owner, repo, &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	})
	return err
}

func (c *githubInstallationClient) PostLabels(ctx context.Context, owner string, repositories []*github.Repository, labels map[string]Label) []error {
	var errors []error

	for _, repository := range repositories {
		for _, label := range labels {
			err := c.PostLabel(ctx, owner, *repository.Name, label)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	return errors
}
