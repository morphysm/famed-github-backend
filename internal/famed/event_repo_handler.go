package famed

import (
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

var labels = []installation.Label{
	{Name: "famed", Color: "566FDB", Description: "Famed - Tracked by Famed"},
	{Name: "none", Color: "566FDB", Description: "Famed - Common Vulnerability Scoring System (CVSS) - None"},
	{Name: "low", Color: "566FDB", Description: "Famed - Common Vulnerability Scoring System (CVSS) - Low"},
	{Name: "medium", Color: "566FDB", Description: "Famed - Common Vulnerability Scoring System (CVSS) - Medium"},
	{Name: "high", Color: "566FDB", Description: "Famed - Common Vulnerability Scoring System (CVSS) - High"},
	{Name: "critical", Color: "566FDB", Description: "Famed - Common Vulnerability Scoring System (CVSS) - Critical"},
}

// handleInstallationRepositoriesEvent adds the labels needed for Famed to the added repository
func (gH *githubHandler) handleInstallationRepositoriesEvent(c echo.Context, event *github.InstallationRepositoriesEvent) error {
	// TODO add event validity check
	if *event.Action != "added" {
		return c.NoContent(http.StatusOK)
	}

	for _, repository := range event.RepositoriesAdded {
		for _, label := range labels {
			err := gH.githubInstallationClient.PostLabel(c.Request().Context(), *event.Installation.Account.Login, *repository.Name, label)
			if err != nil {
				log.Printf("[handleInstallationRepositoriesEvent] error while posting label: %v", err)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
