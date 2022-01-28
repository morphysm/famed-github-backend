package kudo

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrMissingRepoPathParameter = errors.New("missing repo name path parameter")

// GetContributors returns a list of contributors for the kudo board.
func (gH *githubHandler) GetContributors(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	boardGenerator := NewBoardGenerator(gH.kudoConfig, gH.githubInstallationClient, gH.currencyClient, repoName)

	contributors, err := boardGenerator.GetContributors(c.Request().Context())
	if err != nil {
		if errors.Is(err, echo.ErrBadGateway.SetInternal(err)) {
			return err
		}

		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}
