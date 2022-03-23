package famed

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrMissingRepoPathParameter  = errors.New("missing name name path parameter")
	ErrMissingOwnerPathParameter = errors.New("missing owner path parameter")

	ErrAppNotInstalled = errors.New("app not installed")
)

// GetContributors returns a list of contributors for the famed board.
func (gH *githubHandler) GetContributors(c echo.Context) error {
	owner := c.Param("owner")
	if owner == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingOwnerPathParameter)
	}

	repoName := c.Param("repo_name")
	if repoName == "" {
		return echo.ErrBadRequest.SetInternal(ErrMissingRepoPathParameter)
	}

	if installed := gH.githubInstallationClient.CheckInstallation(owner); !installed {
		log.Printf("[Contributors] error on request for contributors: %v", ErrAppNotInstalled)
		return ErrAppNotInstalled
	}

	issues, err := gH.loadIssuesAndEvents(c.Request().Context(), owner, repoName)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return c.JSON(http.StatusOK, []*Contributor{})
	}

	// Use issues with events to generate contributor list
	contributors := contributorsArray(issues, BoardOptions{
		currency: gH.famedConfig.Currency,
		rewards:  gH.famedConfig.Rewards,
	})

	return c.JSON(http.StatusOK, contributors)
}
