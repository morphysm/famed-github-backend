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
// TODO refactor
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	if _, err := kudo.IsValidCloseEvent(event, gH.kudoLabel); err != nil {
		switch err {
		case kudo.ErrIssueMissingAssignee:
			comment := kudo.GenerateCommentFromError(err)
			_, err = gH.githubInstallationClient.PostComment(c.Request().Context(), *event.Repo.Name, *event.Issue.Number, comment)
			if err != nil {
				log.Printf("[handleIssueEvent] error while posting comment: %v", err)
				return err
			}
		default:
			return c.NoContent(http.StatusOK)
		}
	}

	// Get issue events
	events, err := gH.githubInstallationClient.GetIssueEvents(c.Request().Context(), *event.Repo.Name, *event.Issue.Number)
	if err != nil {
		log.Printf("[handleIssueEvent] error getting issue events: %v", err)
		return err
	}

	usdToEthRate, err := gH.currencyClient.GetUSDToETHConversion(c.Request().Context())
	if err != nil {
		log.Printf("[handleIssueEvent] error getting usd eth conversion rate: %v", err)
		return err
	}

	comment := kudo.GenerateComment(event.Issue, events, gH.kudoRewardUnit, gH.kudoRewards, usdToEthRate)

	_, err = gH.githubInstallationClient.PostComment(c.Request().Context(), *event.Repo.Name, *event.Issue.Number, comment)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting comment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
