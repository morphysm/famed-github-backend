package famed

import (
	"context"
	"log"

	"github.com/google/go-github/v41/github"
)

func (gH githubHandler) postLabels(ctx context.Context, repositories []*github.Repository, owner string) {
	for _, repository := range repositories {
		for _, label := range gH.famedConfig.Labels {
			err := gH.githubInstallationClient.PostLabel(ctx, owner, *repository.Name, label)
			if err != nil {
				log.Printf("[handleInstallationRepositoriesEvent] error while posting label: %v", err)
			}
		}
	}
}
