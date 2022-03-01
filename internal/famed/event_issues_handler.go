package famed

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// handleIssuesEvent handles issue events and posts a suggested payout comment to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	ctx := c.Request().Context()

	comment, err := gH.eventToComment(ctx, event)
	if err != nil {
		return err
	}

	// Post comment to GitHub
	err = gH.githubInstallationClient.PostComment(ctx, *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number, comment)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting comment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (gH *githubHandler) eventToComment(ctx context.Context, event *github.IssuesEvent) (string, error) {
	_, err := IsValidCloseEvent(event, gH.famedConfig.Label)
	if err != nil {
		if errors.Is(err, ErrIssueMissingAssignee) {
			return commentFromError(err), nil
		}
		return "", err
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.GetComment(ctx, event.Issue)
}
