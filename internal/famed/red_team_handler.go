package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/config"
	model2 "github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
)

func (gH *githubHandler) GetRedTeam(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model2.ErrMissingOwnerPathParameter.Error())
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model2.ErrMissingRepoPathParameter.Error())
	}

	if ok := gH.githubInstallationClient.CheckInstallation(owner); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, model2.ErrAppNotInstalled.Error())
	}

	famedLabel := gH.famedConfig.Labels[config.FamedLabelKey]
	issueState := model.All
	issues, err := gH.githubInstallationClient.GetIssuesByRepo(c.Request().Context(), owner, repoName, []string{famedLabel.Name}, &issueState)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	redTeam, err := model2.NewRedTeamFromIssues(issues, gH.famedConfig.Currency, gH.now())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, redTeam)
}
