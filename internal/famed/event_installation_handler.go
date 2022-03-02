package famed

import (
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationRepositoriesEvent(c echo.Context, event *github.InstallationRepositoriesEvent) error {
	if valid, err := isValidRepoAddedEvent(event); !valid {
		log.Printf("[handleInstallationRepositoriesEvent] error is not valid repo added event: %v", err)
		return err
	}

	for _, repository := range event.RepositoriesAdded {
		for _, label := range gH.famedConfig.Labels {
			err := gH.githubInstallationClient.PostLabel(c.Request().Context(), *event.Installation.Account.Login, *repository.Name, label)
			if err != nil {
				log.Printf("[handleInstallationRepositoriesEvent] error while posting label: %v", err)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
