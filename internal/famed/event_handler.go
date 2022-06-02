package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

// PostEvent receives the events send to the webhook set in the GitHub App.
// IssueEvents are handled by handleIssuesEvent.
// All other events are ignored.
func (gH *githubHandler) PostEvent(c echo.Context) error {
	event, err := gH.githubInstallationClient.ValidateWebHookEvent(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch event := event.(type) {
	case model.IssuesEvent:
		gH.issuesEventWG.Wait(event.Issue.ID)
		defer gH.issuesEventWG.Done(event.Issue.ID)
		return gH.handleIssuesEvent(c, event)
	case model.InstallationRepositoriesEvent:
		return gH.handleInstallationRepositoriesEvent(c, event)
	case model.InstallationEvent:
		return gH.handleInstallationEvent(c, event)
	default:
		log.Printf("received unhandled event: %v\n", event)
		return c.NoContent(http.StatusOK)
	}
}
