package model

import "github.com/google/go-github/v41/github"

type InstallationRepositoriesEvent struct {
	Action            string
	Installation      RepositoriesInstallation
	RepositoriesAdded []Repository
}

type RepositoriesInstallation struct {
	Account User
}

func NewInstallationRepositoriesEvent(event *github.InstallationRepositoriesEvent) (InstallationRepositoriesEvent, error) {
	if event == nil ||
		event.Action == nil ||
		event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil {
		return InstallationRepositoriesEvent{}, ErrEventMissingData
	}

	account, err := NewUser(event.Installation.Account)
	if err != nil {
		return InstallationRepositoriesEvent{}, err
	}

	compressedEvent := InstallationRepositoriesEvent{
		Action: *event.Action,
		Installation: RepositoriesInstallation{
			Account: account,
		},
	}

	for _, repository := range event.RepositoriesAdded {
		compressedEvent.RepositoriesAdded = append(compressedEvent.RepositoriesAdded, Repository{Name: *repository.Name})
	}

	return compressedEvent, nil
}
