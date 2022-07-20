package github

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/famed/model"
)

func (gH *githubHandler) GetRateLimits(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMissingOwnerPathParameter.Error())
	}

	installations, err := gH.githubInstallationClient.GetRateLimits(c.Request().Context(), owner)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, installations)
}
