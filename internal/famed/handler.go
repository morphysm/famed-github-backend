package famed

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/app"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type HTTPHandler interface {
	GetContributors(c echo.Context) error
	GetInstallations(c echo.Context) error

	PostEvent(c echo.Context) error

	UpdateComments(c echo.Context) error
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubAppClient          app.Client
	githubInstallationClient installation.Client

	famedConfig Config
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubAppClient app.Client, githubInstallationClient installation.Client, famedConfig Config) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
		famedConfig:              famedConfig,
	}
}
