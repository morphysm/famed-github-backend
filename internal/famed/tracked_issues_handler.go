package famed

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/github"
)

// GetTrackedIssues returns all issues tracked.
func (gH *githubHandler) GetTrackedIssues(c echo.Context) error {
	ctx := c.Request().Context()

	installations, err := gH.githubAppClient.GetInstallations(ctx)
	if err != nil {
		return err
	}

	trackedIssues := make(map[string][]github.Issue)
	for _, installation := range installations {
		repos, err := gH.githubInstallationClient.GetRepos(ctx, installation.Account.Login)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			issueState := github.All
			issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, installation.Account.Login, repo.Name, []string{"famed"}, &issueState)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%s/%s", installation.Account.Login, repo.Name)
			if issues == nil {
				// add empty slice so that the json response contains [] instead of nil
				trackedIssues[key] = []github.Issue{}
				continue
			}
			trackedIssues[key] = issues
		}
	}

	return c.JSON(http.StatusOK, trackedIssues)
}
