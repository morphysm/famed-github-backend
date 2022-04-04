package famed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (gH *githubHandler) GetRedTeam(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingOwnerPathParameter)
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	redTeam, err := readRedTeamFromFile(owner, repoName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, redTeam)
}

func readRedTeamFromFile(owner string, repo string) ([]*Contributor, error) {
	path := fmt.Sprintf("redTeams/%s/%s/redTeam.json", owner, repo)

	rawRedTeam, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	redTeam := Contributors{}
	err = json.Unmarshal(rawRedTeam, &redTeam)
	if err != nil {
		return nil, err
	}

	redTeam.updateMonthlyRewards()
	return redTeam.toSortedSlice(), nil
}
