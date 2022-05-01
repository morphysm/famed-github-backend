package model

import "github.com/google/go-github/v41/github"

type InstallationEvent struct {
	Action       string
	Repositories []Repository
	Installation
}

func NewInstallationEvent(event *github.InstallationEvent) (InstallationEvent, error) {
	if event == nil ||
		event.Action == nil ||
		event.Installation == nil ||
		event.Installation.Account == nil ||
		event.Installation.Account.Login == nil ||
		event.Installation.ID == nil {
		return InstallationEvent{}, ErrEventMissingData
	}

	account, err := NewUser(event.Installation.Account)
	if err != nil {
		return InstallationEvent{}, err
	}

	compressedEvent := InstallationEvent{
		Action: *event.Action,
		Installation: Installation{
			ID:      *event.Installation.ID,
			Account: account,
		},
	}

	for _, repository := range event.Repositories {
		compressedEvent.Repositories = append(compressedEvent.Repositories, Repository{Name: *repository.Name})
	}

	return compressedEvent, nil
}
