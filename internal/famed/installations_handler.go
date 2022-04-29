package famed

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/client/github"
)

type installation struct {
	github.Installation
	Repos []string `json:"repositories"`
}

func (gH *githubHandler) GetInstallations(c echo.Context) error {
	installations, err := gH.githubAppClient.GetInstallations(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	resp := make([]installation, len(installations))
	for i, instal := range installations {
		repositories, err := gH.githubInstallationClient.GetRepos(c.Request().Context(), instal.Account.Login)
		if err != nil {
			return err
		}

		repos := make([]string, len(repositories))
		for j, repo := range repositories {
			repos[j] = repo.Name
		}

		resp[i] = installation{Installation: instal, Repos: repos}
	}

	return c.JSON(http.StatusOK, resp)
}
