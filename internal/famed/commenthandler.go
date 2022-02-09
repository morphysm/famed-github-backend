package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// UpdateComments updates the comments in a GitHub name.
func (gH *githubHandler) UpdateComments(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, repoName)

	comments, err := repo.GetComments(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, comments)
}
