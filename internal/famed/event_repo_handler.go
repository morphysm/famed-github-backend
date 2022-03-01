package famed

import (
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var labels = []installation.Label{
	{Name: "famed", Color: "566FDB", Description: "Issues with this label will be tracked by Famed"},
	{Name: "none", Color: "566FDB", Description: "Common Vulnerability Scoring System (CVSS) label to be used with Famed"},
	{Name: "low", Color: "566FDB", Description: "Common Vulnerability Scoring System (CVSS) label to be used with Famed"},
	{Name: "medium", Color: "566FDB", Description: "Common Vulnerability Scoring System (CVSS) label to be used with Famed"},
	{Name: "high", Color: "566FDB", Description: "Common Vulnerability Scoring System (CVSS) label to be used with Famed"},
	{Name: "critical", Color: "566FDB", Description: "Common Vulnerability Scoring System (CVSS) label to be used with Famed"},
}

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationRepositoriesEvent(c echo.Context, event *github.InstallationRepositoriesEvent) error {
	// TODO add event validity check
	if *event.Action != "added" {
		return c.NoContent(http.StatusOK)
	}

	for _, repository := range event.RepositoriesAdded {
		for _, label := range labels {
			gH.githubInstallationClient.PostLabel(c.Request().Context(), *repository.Owner.Login, *repository.Name, label)
		}
	}

	return c.NoContent(http.StatusOK)
}
