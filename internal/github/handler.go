package github

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
)

type HTTPHandler interface {
	GetRateLimit(c echo.Context) error
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubInstallationClient providers.InstallationClient
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubInstallationClient providers.InstallationClient) HTTPHandler {
	return &githubHandler{
		githubInstallationClient: githubInstallationClient,
	}
}
