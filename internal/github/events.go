package github

import (
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

// PostEvent receives the events send to the webhook set in the GitHub App.
// IssueEvents are handled by handleIssuesEvent.
// All other events are ignored.
func (gH *githubHandler) PostEvent(c echo.Context) error {
	payload, err := github.ValidatePayload(c.Request(), []byte(gH.webhookSecret))
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
// if the kudo label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	generator := kudo.NewCommentGenerator(gH.kudoConfig, gH.githubInstallationClient, gH.currencyClient, event)

	comment, err := generator.GetComment(c.Request().Context())
	if err != nil {
		log.Printf("[handleIssueEvent] error while generating comment: %v", err)
		return err
	}

	// Post comment to GitHub
	_, err = gH.githubInstallationClient.PostComment(c.Request().Context(), *event.Repo.Name, *event.Issue.Number, comment)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting comment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
