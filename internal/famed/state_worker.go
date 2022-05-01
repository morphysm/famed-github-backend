package famed

import (
	"context"
	"log"
)

// CleanState iterates over all issues and updates their comments if necessary.
func (gH *githubHandler) CleanState() {
	log.Printf("[CleanState] running clean up...")

	ctx := context.Background()
	installations, err := gH.githubAppClient.GetInstallations(ctx)
	if err != nil {
		log.Printf("[CleanState] error while getting installations: %v", err)
	}

	for _, installation := range installations {
		// Check if installation client is set up and if necessary add client
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

		for _, repoName := range repos {
			go func(owner string, repoName string) {
				err := gH.updateRewardComments(ctx, owner, repoName, nil)
				if err != nil {
					log.Printf("[CleanState] error while updating reward comments: %v", err)
				}
			}(installation.Account.Login, repoName)
			go func(owner string, repoName string) {
				err := gH.updateEligibleComments(ctx, owner, repoName, nil)
				if err != nil {
					log.Printf("[CleanState] error while updating reward comments: %v", err)
				}
			}(installation.Account.Login, repoName)
		}
	}
}
