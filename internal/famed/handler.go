package famed

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/currency"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type HTTPHandler interface {
	GetContributors(c echo.Context) error
	PostEvent(c echo.Context) error
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubInstallationClient installation.Client
	currencyClient           currency.Client
	webhookSecret            *string
	installationID           int64
	famedConfig              Config
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubInstallationClient installation.Client, currencyClient currency.Client, webhookSecret *string, installationID int64, famedConfig Config) HTTPHandler {
	return &githubHandler{
		githubInstallationClient: githubInstallationClient,
		currencyClient:           currencyClient,
		webhookSecret:            webhookSecret,
		installationID:           installationID,
		famedConfig:              famedConfig,
	}
}
