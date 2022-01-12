package github

import (
	"errors"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

func (gH *githubHandler) GetEvents(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo id path parameter"))
	}

	eventsResp, err := gH.githubInstallationClient.GetRepoEvents(c.Request().Context(), repoName)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, eventsResp)
}

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
		log.Debugf("received unhandled event: %v", event)

		return c.NoContent(http.StatusOK)
	}
}

// handleIssuesEvent handles issue events.
// If the kudo label is set and the issue is closed a suggested payout comment is posted to the GitHub API.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	if event.Action == nil ||
		*event.Action != string(installation.Closed) ||
		event.Repo.Name == nil ||
		event.Issue.Number == nil {
		return c.NoContent(http.StatusOK)
	}

	// TODO Check for labels and alike
	repoName := *event.Repo.Name
	issueNumber := *event.Issue.Number

	testComment := "This will be suggested payout"

	_, err := gH.githubInstallationClient.PostComment(c.Request().Context(), repoName, issueNumber, testComment)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
