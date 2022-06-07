package model

import "github.com/google/go-github/v41/github"

func NewRepo(repo *github.Repository) (string, error) {
	if repo == nil ||
		repo.Name == nil {
		return "", ErrRepoMissingData
	}

	return *repo.Name, nil
}
