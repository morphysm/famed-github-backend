package github

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/go-github/v41/github"
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client/installation"
)

func (gH *githubHandler) GetIssues(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo id path parameter"))
	}

	issuesResp, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), repoName, []string{"payed"}, installation.Closed)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, issuesResp)
}

func (gH *githubHandler) PostComment(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo id path parameter"))
	}

	issueNumber, err := strconv.Atoi(c.Param("issue_number"))
	if err != nil {
		return echo.ErrBadRequest.SetInternal(errors.New("missing or incorrect issue number path parameter"))
	}

	var comment github.IssueComment
	err = c.Bind(&comment)
	if err != nil {
		return err
	}

	if comment.Body == nil {
		return echo.ErrBadRequest.SetInternal(errors.New("missing comment body"))
	}

	commentResponse, err := gH.githubInstallationClient.PostComment(c.Request().Context(), repoName, issueNumber, *comment.Body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, commentResponse)
}