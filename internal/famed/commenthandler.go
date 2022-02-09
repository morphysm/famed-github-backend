package famed

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UpdateComments updates the comments in a GitHub name.
func (gH *githubHandler) UpdateComments(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	boardGenerator := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, repoName)

	contributors, err := boardGenerator.GetContributors(c.Request().Context())
	if err != nil {
		if errors.Is(err, echo.ErrBadGateway) {
			return err
		}

		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}
