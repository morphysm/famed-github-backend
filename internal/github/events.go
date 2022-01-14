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

// handleIssuesEvent handles issue events.
// If the kudo label is set and the issue is closed a suggested payout comment is posted to the GitHub API.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	if !kudo.IsValidCloseEvent(event, gH.kudoLabel) {
		return c.NoContent(http.StatusOK)
	}

	// Get issue events
	events, err := gH.githubInstallationClient.GetIssueEvents(c.Request().Context(), *event.Repo.Name, *event.Issue.Number)
	if err != nil {
		log.Printf("error getting issue events: %v", err)
		return err
	}

	contributors := kudo.GenerateContributorsByIssue(nil, event.Issue, events)
	comment := kudo.GenerateCommentFromContributors(contributors)

	_, err = gH.githubInstallationClient.PostComment(c.Request().Context(), *event.Repo.Name, *event.Issue.Number, comment)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
