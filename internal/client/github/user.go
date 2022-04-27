package github

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

type User struct {
	Login     string
	AvatarURL string
	HTMLURL   string
}

// getUser gets a GitHub user.
func (c *githubInstallationClient) GetUser(ctx context.Context, owner string, login string) (User, error) {
	client, _ := c.clients.get(owner)

	user, _, err := client.Users.Get(ctx, login)
	if err != nil {
		return User{}, err
	}

	validUser, err := validateUser(user)
	if err != nil {
		return User{}, err
	}

	return validUser, nil
}

func validateUser(user *github.User) (User, error) {
	if user == nil ||
		user.Login == nil {
		return User{}, ErrUserMissingData
	}

	return User{
		Login:     *user.Login,
		AvatarURL: pointer.ToString(user.AvatarURL),
		HTMLURL:   pointer.ToString(user.HTMLURL),
	}, nil
}
