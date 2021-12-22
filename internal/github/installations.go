package github

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetInstallations(c echo.Context) error {
	installations, err := gH.githubClient.GetInstallations(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, installations)
}
