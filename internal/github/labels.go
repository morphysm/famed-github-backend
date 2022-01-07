package github

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetLabels(c echo.Context) error {
	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo id path parameter"))
	}

	labelResp, err := gH.githubInstallationClient.GetLabels(c.Request().Context(), repoName)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, labelResp)
}
