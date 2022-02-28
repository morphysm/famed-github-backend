package famed

import (
	"errors"
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
	default:
		log.Printf("received unhandled event: %v\n", event)
		return c.NoContent(http.StatusOK)
	}
}

// handleIssuesEvent handles issue events and posts a suggested payout comment to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	comment, err := gH.eventToComment(c, event)
	if err != nil {
		return err
	}

	// Post comment to GitHub
	_, err = gH.githubInstallationClient.PostComment(c.Request().Context(), *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number, comment)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting comment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (gH *githubHandler) eventToComment(c echo.Context, event *github.IssuesEvent) (string, error) {
	_, err := IsValidCloseEvent(event, gH.famedConfig.Label)
	if err != nil {
		if errors.Is(err, ErrIssueMissingAssignee) {
			return commentFromError(err), nil
		}
		return "", err
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.GetComment(c.Request().Context(), event.Issue)
}
