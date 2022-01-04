package github

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAccessTokens returns access tokens for all repos accessible by an installation
func (gH *githubHandler) GetAccessTokens(c echo.Context) error {
	installationID := c.Param("installation_id")
	if installationID == "" {
		return echo.ErrBadRequest.SetInternal(errors.New("missing installation id path parameter"))
	}

	// Get installation token without repo access
	accessTokensResp, err := gH.githubClient.GetAccessTokens(c.Request().Context(), installationID, nil)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	// Get repos
	reposResp, err := gH.githubClient.GetRepos(c.Request().Context(), accessTokensResp.Token)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	// Generate list of repo ids
	var repoIDs []int
	for _, repository := range reposResp.Repositories {
		repoIDs = append(repoIDs, repository.Id)
	}

	// Get installation token with repo access
	accessTokensResp, err = gH.githubClient.GetAccessTokens(c.Request().Context(), installationID, repoIDs)
	if err != nil {
		return echo.ErrBadGateway.SetInternal(err)
	}

	return c.JSON(http.StatusOK, accessTokensResp)
}
