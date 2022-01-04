package github

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/app"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

type HTTPHandler interface {
	GetInstallations(c echo.Context) error
	GetLabels(c echo.Context) error
	GetRepos(c echo.Context) error
}

// githubHandler represents the handler for the github endpoints.
type githubHandler struct {
	githubAppClient          app.Client
	githubInstallationClient installation.Client
}

// NewHandler returns a pointer to the github handler.
func NewHandler(githubAppClient app.Client, githubInstallationClient installation.Client) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
	}
}
