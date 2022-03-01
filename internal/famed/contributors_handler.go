package famed

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrMissingRepoPathParameter  = errors.New("missing name name path parameter")
	ErrMissingOwnerPathParameter = errors.New("missing owner path parameter")
)

// GetContributors returns a list of contributors for the famed board.
func (gH *githubHandler) GetContributors(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingOwnerPathParameter)
	}
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	repo := NewRepo(gH.famedConfig, gH.githubInstallationClient, gH.currencyClient, owner, repoName)

	contributors, err := repo.GetContributors(c.Request().Context())
	if err != nil {
		if errors.Is(err, echo.ErrBadGateway) {
			return err
		}

		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, contributors)
}
