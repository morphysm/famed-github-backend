package famed

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationEvent(c echo.Context, event installation.InstallationEvent) error {
	if event.Action != "created" {
		log.Printf("[handleInstallationEvent] error is not valid insatllation created event")
		return ErrEventNotInstallationCreated
	}

	err := gH.githubInstallationClient.AddInstallation(event.Installation.Account.Login, event.Installation.ID)
	if err != nil {
		return err
	}

	var repoNames []string
	for _, repository := range event.Repositories {
		repoNames = append(repoNames, repository.Name)
	}
	errors := gH.githubInstallationClient.PostLabels(c.Request().Context(), event.Installation.Account.Login, repoNames, gH.famedConfig.Labels)
	for _, err := range errors {
		log.Printf("[handleInstallationEvent] error while posting labels: %v", err)
	}

	return c.NoContent(http.StatusOK)
}
