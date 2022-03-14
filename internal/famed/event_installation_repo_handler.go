package famed

import (
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationRepositoriesEvent(c echo.Context, event *github.InstallationRepositoriesEvent) error {
	if valid, err := isRepoAddedEventValid(event); !valid {
		log.Printf("[handleInstallationRepositoriesEvent] error is not valid repo added event: %v", err)
		return err
	}

	errors := gH.githubInstallationClient.PostLabels(c.Request().Context(), *event.Installation.Account.Login, event.RepositoriesAdded, gH.famedConfig.Labels)
	for _, err := range errors {
		log.Printf("[handleInstallationRepositoriesEvent] error while posting labels: %v", err)
	}

	return c.NoContent(http.StatusOK)
}
