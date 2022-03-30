package famed

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/github"
)

type HTTPHandler interface {
	GetInstallations(c echo.Context) error
	GetRateLimit(c echo.Context) error

	GetContributors(c echo.Context) error
	GetRedTeam(c echo.Context) error

	PostEvent(c echo.Context) error

	GetUpdateComments(c echo.Context) error

	CleanState()
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubAppClient          github.AppClient
	githubInstallationClient github.InstallationClient

	famedConfig Config
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubAppClient github.AppClient, githubInstallationClient github.InstallationClient, famedConfig Config) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
		famedConfig:              famedConfig,
	}
}
