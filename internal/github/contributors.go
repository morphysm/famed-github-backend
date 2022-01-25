package github

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

// GetContributors returns a list of contributors for the kudo board.
func (gH *githubHandler) GetContributors(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo name path parameter"))
	}

	boardGenerator := kudo.NewBoardGenerator(gH.kudoConfig, gH.githubInstallationClient, gH.currencyClient, repoName)

	contributors, err := boardGenerator.GetContributors(c.Request().Context())
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}
