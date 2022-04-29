package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetRateLimit(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, ErrMissingOwnerPathParameter.Error())
	}

	installations, err := gH.githubInstallationClient.GetRateLimit(c.Request().Context(), owner)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, installations)
}
