package famed

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type HTTPHandler interface {
	GetContributors(c echo.Context) error
	PostEvent(c echo.Context) error
	UpdateComments(c echo.Context) error
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubInstallationClient installation.Client
	webhookSecret            *string
	famedConfig              Config
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubInstallationClient installation.Client, webhookSecret *string, famedConfig Config) HTTPHandler {
	return &githubHandler{
		githubInstallationClient: githubInstallationClient,
		webhookSecret:            webhookSecret,
		famedConfig:              famedConfig,
	}
}
