package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"

	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationEvent(c echo.Context, event model.InstallationEvent) error {
	if event.Action != "created" {
		log.Error().Msg("[handleInstallationEvent] error is not valid installation created event")
		return model2.ErrEventNotInstallationCreated
	}

	err := gH.githubInstallationClient.AddInstallation(event.Installation.Account.Login, event.Installation.ID)
	if err != nil {
		return err
	}

	repoNames := make([]string, len(event.Repositories))
	for i, repository := range event.Repositories {
		repoNames[i] = repository.Name
	}
	errors := gH.githubInstallationClient.PostLabels(c.Request().Context(), event.Installation.Account.Login, repoNames, gH.famedConfig.Labels)
	for _, err := range errors {
		log.Error().Err(err).Msg("[handleInstallationEvent] error while posting labels")
	}

	return c.NoContent(http.StatusOK)
}
