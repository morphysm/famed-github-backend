package github

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
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
	if event.Action == nil ||
		*event.Action != string(installation.Closed) ||
		event.Repo == nil ||
		event.Issue == nil ||
		event.Repo.Name == nil ||
		event.Issue.Number == nil ||
		event.Issue.Labels == nil {
		return c.NoContent(http.StatusOK)
	}

	// TODO Check for labels and alike
	repoName := *event.Repo.Name
	issueNumber := *event.Issue.Number

	// Check labels for "kudo" and severity
	kudoSupported := false
	for _, label := range event.Issue.Labels {
		if label.Name != nil && *label.Name == gH.kudoLabel {
			kudoSupported = true
		}
	}

	if !kudoSupported {
		return c.NoContent(http.StatusOK)
	}

	severity := kudo.IssueToSeverity(event.Issue)

	// Get issue events
	events, err := gH.githubInstallationClient.GetIssueEvents(c.Request().Context(), repoName, issueNumber)
	if err != nil {
		log.Printf("error getting issue events: %v", err)
		return err
	}

	// TODO not optimal to reuse the function for all contributors of a repo
	contributors := kudo.EventsToContributors(nil, events, *event.Issue.CreatedAt, *event.Issue.ClosedAt, severity)
	comment := "Kudo suggests:"
	for _, contributor := range contributors {
		comment = fmt.Sprintf("%s\n Contributor: %s, Reward: %f\n", comment, contributor.Login, contributor.RewardSum)
	}

	_, err = gH.githubInstallationClient.PostComment(c.Request().Context(), repoName, issueNumber, comment)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
