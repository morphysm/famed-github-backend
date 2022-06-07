package model

import "github.com/google/go-github/v41/github"

type Installation struct {
	ID      int64
	Account User
}

func NewInstallation(installation *github.Installation) (Installation, error) {
	if installation == nil ||
		installation.ID == nil {
		return Installation{}, ErrInstallationMissingData
	}

	account, err := NewUser(installation.Account)
	if err != nil {
		return Installation{}, err
	}

	return Installation{ID: *installation.ID, Account: account}, nil
}
