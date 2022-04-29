package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetContributors returns a list of contributors for the famed board.
func (gH *githubHandler) GetContributors(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, ErrMissingOwnerPathParameter.Error())
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, ErrMissingRepoPathParameter.Error())
	}

	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, ErrAppNotInstalled.Error())
	}

	issues, err := gH.loadIssues(c.Request().Context(), owner, repoName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	if len(issues) == 0 {
		return c.JSON(http.StatusOK, []*contributor{})
	}

	// Use issues with events to generate contributor list
	contributors := contributorsArray(issues, boardOptions{
		currency:  gH.famedConfig.Currency,
		rewards:   gH.famedConfig.Rewards,
		daysToFix: gH.famedConfig.DaysToFix,
	})

	return c.JSON(http.StatusOK, contributors)
}
