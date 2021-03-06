package providers

import (
	"context"

	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

func (c *githubInstallationClient) PostLabel(ctx context.Context, owner string, repoName string, label model.Label) error {
	client, _ := c.clients.get(owner)

	_, _, err := client.Issues.CreateLabel(ctx, owner, repoName, &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	})
	return err
}

func (c *githubInstallationClient) PostLabels(ctx context.Context, owner string, repoNames []string, labels map[string]model.Label) []error {
	var errors []error

	for _, repo := range repoNames {
		for _, label := range labels {
			err := c.PostLabel(ctx, owner, repo, label)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	return errors
}
