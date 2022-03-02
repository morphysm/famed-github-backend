package famed

import (
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationEvent(c echo.Context, event *github.InstallationEvent) error {
	if valid, err := isValidInstallationEvent(event); !valid {
		log.Printf("[handleInstallationEvent] error is not valid insatllation event: %v", err)
		return err
	}

	err := gH.githubInstallationClient.AddInstallation(*event.Installation.Account.Login, *event.Installation.ID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
