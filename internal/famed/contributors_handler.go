package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/famed/model"
)

// GetBlueTeam returns a list of contributors for the famed board.
func (gH *githubHandler) GetBlueTeam(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMissingOwnerPathParameter.Error())
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMissingRepoPathParameter.Error())
	}

	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrAppNotInstalled.Error())
	}

	issues, err := gH.loadEnrichedIssues(c.Request().Context(), owner, repoName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	if len(issues) == 0 {
		return c.JSON(http.StatusOK, []*model.Contributor{})
	}

	// Use issues with events to generate contributor list
	rewardStructure := model.NewRewardStructure(gH.famedConfig.Rewards, gH.famedConfig.DaysToFix, 2)
	boardOptions := model.NewBoardOptions(gH.famedConfig.Currency, rewardStructure, gH.now())
	contributors := model.NewBlueTeamFromIssues(issues, boardOptions)

	return c.JSON(http.StatusOK, contributors)
}
