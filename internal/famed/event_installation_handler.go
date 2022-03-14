package famed

import (
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationEvent(c echo.Context, event *github.InstallationEvent) error {
	if valid, err := isInstallationEventValid(event); !valid {
		log.Printf("[handleInstallationEvent] error is not valid insatllation created event: %v", err)
		return err
	}

	err := gH.githubInstallationClient.AddInstallation(*event.Installation.Account.Login, *event.Installation.ID)
	if err != nil {
		return err
	}

	errors := gH.githubInstallationClient.PostLabels(c.Request().Context(), *event.Installation.Account.Login, event.Repositories, gH.famedConfig.Labels)
	for _, err := range errors {
		log.Printf("[handleInstallationEvent] error while posting labels: %v", err)
	}

	return c.NoContent(http.StatusOK)
}
