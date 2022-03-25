package famed

import (
	"context"
	"log"
)

func (gH *githubHandler) CleanState() {
	log.Printf("[CleanState] running clean up...")

	ctx := context.Background()
	installations, err := gH.githubAppClient.GetInstallations(ctx)
	if err != nil {
		log.Printf("[CleanState] error while getting installations: %v", err)
	}

	for _, installation := range installations {
		// Check and if necessary add github clients
		//TODO add check for null pointer
		if !gH.githubInstallationClient.CheckInstallation(installation.Account.Login) {
			err := gH.githubInstallationClient.AddInstallation(installation.Account.Login, installation.ID)
			if err != nil {
				log.Printf("[CleanState] error while adding github: %v", err)
				continue
			}
		}

		repos, err := gH.githubInstallationClient.GetRepos(ctx, installation.Account.Login)
		if err != nil {
			log.Printf("[CleanState] error while getting repos: %v", err)
			continue
		}

		for _, repo := range repos {
			go gH.updateRewardComments(ctx, installation.Account.Login, repo.Name, nil)
			go gH.updateEligibleComments(ctx, installation.Account.Login, repo.Name, nil)
		}
	}
}
