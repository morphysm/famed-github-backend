package famed

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
)

// handleIssuesEvent handles issue events and posts a suggested payout comment to the GitHub API,
// if the famed label is set and the issue is closed.
func (gH *githubHandler) handleIssuesEvent(c echo.Context, event *github.IssuesEvent) error {
	ctx := c.Request().Context()

	var comment string
	var err error
	//TODO add action check
	switch *event.Action {
	case string(installation.Closed):
		comment, err = gH.closeEventToComment(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating comment for closed event: %v", err)
			// TODO return c.NoContent(http.StatusOK)
			return err
		}
	case string(installation.Labeled):
		comment, err = gH.labeledEventToComment(ctx, event)
		if err != nil {
			log.Printf("[handleIssuesEvent] error while generating comment for labeled event: %v", err)
			// TODO return c.NoContent(http.StatusOK)
			return err
		}
	default:
		log.Print("received unhandled issues event")
		return c.NoContent(http.StatusOK)
	}

	// Post comment to GitHub
	err = gH.githubInstallationClient.PostComment(ctx, *event.Repo.Owner.Login, *event.Repo.Name, *event.Issue.Number, comment)
	if err != nil {
		log.Printf("[handleIssueEvent] error while posting comment: %v", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (gH *githubHandler) closeEventToComment(ctx context.Context, event *github.IssuesEvent) (string, error) {
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]

	_, err := isValidCloseEvent(event, famedLabel.Name)
	if err != nil {
		if errors.Is(err, ErrIssueMissingAssignee) {
			return commentFromError(err), nil
		}
		return "", err
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.ContributorComment(ctx, event.Issue)
}

func (gH *githubHandler) labeledEventToComment(ctx context.Context, event *github.IssuesEvent) (string, error) {
	famedLabel := gH.famedConfig.Labels[config.FamedLabel]

	// TODO add severity labels
	// TODO add validity check
	if event.Label == nil || event.Label.Name == nil || *event.Label.Name != famedLabel.Name {
		return "", errors.New("label is not \"famed\" label")
	}

	// TODO move this
	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, *event.Repo.Owner.Login, *event.Repo.Name)
	return repo.IssueStateComment(ctx, event.Issue)
}
