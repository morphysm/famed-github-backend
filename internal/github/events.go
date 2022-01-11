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

type WebhookEvent struct {
	Action *string       `json:"action"`
	Issue  *github.Issue `json:"issue"`

	// TODO look into label changes
	Changes *struct {
	} `json:"changes"`
	Repository *github.Repository `json:"repository"`
	Sender     *github.User       `json:"sender"`
}

func (gH *githubHandler) PostEvent(c echo.Context) error {
	var webhookEvent WebhookEvent

	if err := c.Bind(&webhookEvent); err != nil {
		log.Debugf("error binding webhook event: %v", err)
	}

	if webhookEvent.Action == nil ||
		*webhookEvent.Action != string(installation.Closed) ||
		webhookEvent.Repository.Name == nil ||
		webhookEvent.Issue.Number == nil {
		return c.NoContent(http.StatusOK)
	}

	repoName := *webhookEvent.Repository.Name
	issueNumber := *webhookEvent.Issue.Number

	testComment := "This will be suggested payout"

	gH.githubInstallationClient.PostComment(c.Request().Context(), repoName, issueNumber, testComment)

	return c.NoContent(http.StatusOK)
}
