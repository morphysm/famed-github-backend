package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"

	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationRepositoriesEvent(c echo.Context, event model.InstallationRepositoriesEvent) error {
	if event.Action != "added" {
		log.Error().Msg("[handleInstallationRepositoriesEvent] error is not valid repo added event")
		return model2.ErrEventNotRepoAdded
	}

	repoNames := make([]string, len(event.RepositoriesAdded))
	for i, repository := range event.RepositoriesAdded {
		repoNames[i] = repository.Name
	}
	errors := gH.githubInstallationClient.PostLabels(c.Request().Context(), event.Installation.Account.Login, repoNames, gH.famedConfig.Labels)
	for _, err := range errors {
		log.Error().Err(err).Msg("[handleInstallationRepositoriesEvent] error while posting labels")
	}

	return c.NoContent(http.StatusOK)
}
