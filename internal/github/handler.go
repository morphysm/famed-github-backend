package github

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/apps"
	"github.com/morphysm/kudos-github-backend/internal/client/currency"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

type HTTPHandler interface {
	GetContributors(c echo.Context) error
	PostEvent(c echo.Context) error
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubAppClient          apps.Client
	githubInstallationClient installation.Client
	currencyClient           currency.Client
	webhookSecret            string
	installationID           int64
	kudoLabel                string
	kudoRewardCurrency       string
	kudoRewards              map[kudo.IssueSeverity]float64
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubAppClient apps.Client, githubInstallationClient installation.Client, currencyClient currency.Client, webhookSecret string, installationID int64, kudoLabel string, kudoRewardCurrency string, kudoRewards map[kudo.IssueSeverity]float64) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
		currencyClient:           currencyClient,
		webhookSecret:            webhookSecret,
		installationID:           installationID,
		kudoLabel:                kudoLabel,
		kudoRewardCurrency:       kudoRewardCurrency,
		kudoRewards:              kudoRewards,
	}
}
