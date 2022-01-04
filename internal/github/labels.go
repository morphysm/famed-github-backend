package github

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetLabels(c echo.Context) error {
	repoID := c.Param("repo_id")
	if repoID == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing repo id path parameter"))
	}

	labelResp, err := gH.githubInstallationClient.GetLabels(c.Request().Context(), repoID)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, labelResp)
}
