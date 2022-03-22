package famed

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationRepositoriesEvent(c echo.Context, event installation.InstallationRepositoriesEvent) error {
	if event.Action != "added" {
		log.Printf("[handleInstallationRepositoriesEvent] error is not valid repo added event")
		return ErrEventNotRepoAdded
	}

	errors := gH.githubInstallationClient.PostLabels(c.Request().Context(), event.Installation.Account.Login, event.RepositoriesAdded, gH.famedConfig.Labels)
	for _, err := range errors {
		log.Printf("[handleInstallationRepositoriesEvent] error while posting labels: %v", err)
	}

	return c.NoContent(http.StatusOK)
}
