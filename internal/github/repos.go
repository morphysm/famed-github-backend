package github

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetRepos(c echo.Context) error {
	repoResp, err := gH.githubInstallationClient.GetRepos(c.Request().Context())
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, repoResp)
}
