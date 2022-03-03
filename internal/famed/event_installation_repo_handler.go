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

	gH.postLabels(c.Request().Context(), event.RepositoriesAdded, *event.Installation.Account.Login)

	return c.NoContent(http.StatusOK)
}
