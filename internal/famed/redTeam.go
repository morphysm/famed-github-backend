package famed

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetRedTeam(c echo.Context) error {
	rawRedTeam, err := os.ReadFile("redTeam.json")
	if err != nil {
		return err
	}

	redTeam := Contributors{}
	err = json.Unmarshal(rawRedTeam, &redTeam)
	if err != nil {
		return err
	}

	redTeam.updateMonthlyRewards()
	sortedRedTeam := redTeam.toSortedSlice()

	return c.JSON(http.StatusOK, sortedRedTeam)
}
