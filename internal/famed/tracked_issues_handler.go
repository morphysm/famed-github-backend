package famed

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

type trackedIssue struct {
	model.Issue
	Comments []model.IssueComment
}

// GetTrackedIssues returns all issues tracked.
func (gH *githubHandler) GetTrackedIssues(c echo.Context) error {
	ctx := c.Request().Context()

	installations, err := gH.githubAppClient.GetInstallations(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	trackedIssues := make(map[string][]trackedIssue)
	for _, installation := range installations {
		repos, err := gH.githubInstallationClient.GetRepos(ctx, installation.Account.Login)
		if err != nil {
			return err
		}

		for _, repoName := range repos {
			issueState := model.All
			issues, err := gH.githubInstallationClient.GetIssuesByRepo(ctx, installation.Account.Login, repoName, []string{"famed"}, &issueState)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%s/%s", installation.Account.Login, repoName)
			if issues == nil {
				// add empty slice so that the json response contains [] instead of nil
				trackedIssues[key] = []trackedIssue{}
				continue
			}
			var commentIssues []trackedIssue
			for _, issue := range issues {
				trackedIssue := trackedIssue{Issue: issue}
				comments, err := gH.githubInstallationClient.GetComments(ctx, installation.Account.Login, repoName, issue.Number)
				if err != nil {
					return err
				}
				for _, comment := range comments {
					if comment.User.Login == gH.famedConfig.BotLogin {
						trackedIssue.Comments = append(trackedIssue.Comments, comment)
					}
				}
				commentIssues = append(commentIssues, trackedIssue)
			}
			trackedIssues[key] = commentIssues
		}
	}

	return c.JSON(http.StatusOK, trackedIssues)
}
