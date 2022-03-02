package famed

import (
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// PostEvent receives the events send to the webhook set in the GitHub App.
// IssueEvents are handled by handleIssuesEvent.
// All other events are ignored.
func (gH *githubHandler) PostEvent(c echo.Context) error {
	var webhookSecret []byte
	if gH.webhookSecret != nil {
		webhookSecret = []byte(*gH.webhookSecret)
	}

	payload, err := github.ValidatePayload(c.Request(), webhookSecret)
	if err != nil {
		return err
	}

	event, err := github.ParseWebHook(github.WebHookType(c.Request()), payload)
	if err != nil {
		return err
	}

	switch event := event.(type) {
	case *github.IssuesEvent:
		return gH.handleIssuesEvent(c, event)
	case *github.InstallationRepositoriesEvent:
		return gH.handleInstallationRepositoriesEvent(c, event)
	case *github.InstallationEvent:
		return gH.handleInstallationEvent(c, event)
	default:
		log.Printf("received unhandled event: %v\n", event)
		return c.NoContent(http.StatusOK)
	}
}
