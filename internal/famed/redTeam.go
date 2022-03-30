package famed

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetRedTeam(c echo.Context) error {
	redTeam, err := os.ReadFile("redTeam.json")
	if err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, redTeam)
}
