package model

import (
	"github.com/google/go-github/v41/github"

	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

type User struct {
	Login     string
	AvatarURL string
	HTMLURL   string
}

func NewUser(user *github.User) (User, error) {
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
