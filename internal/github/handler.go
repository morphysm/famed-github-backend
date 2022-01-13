package github

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/apps"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

type HTTPHandler interface {
	GetContributors(c echo.Context) error

	PostEvent(c echo.Context) error
}

// githubHandler represents the handler for the github endpoints.
type githubHandler struct {
	githubAppClient          apps.Client
	githubInstallationClient installation.Client
	webhookSecret            string
	installationID           int64
	kudoLabel                string
}

// NewHandler returns a pointer to the github handler.
func NewHandler(githubAppClient apps.Client, githubInstallationClient installation.Client, webhookSecret string, installationID int64, kudoLabel string) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
		webhookSecret:            webhookSecret,
		installationID:           installationID,
		kudoLabel:                kudoLabel,
	}
}
