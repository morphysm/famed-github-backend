package providers

import (
	"context"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// GetUser returns a GitHub user for a given login.
func (c *githubInstallationClient) GetUser(ctx context.Context, owner string, login string) (model.User, error) {
	client, _ := c.clients.get(owner)

	user, _, err := client.Users.Get(ctx, login)
	if err != nil {
		return model.User{}, err
	}

	validUser, err := model.NewUser(user)
	if err != nil {
		return model.User{}, err
	}

	return validUser, nil
}
