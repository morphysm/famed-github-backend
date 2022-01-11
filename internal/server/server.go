package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/morphysm/kudos-github-backend/internal/client/apps"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/config"
	glib "github.com/morphysm/kudos-github-backend/internal/github"
	"github.com/morphysm/kudos-github-backend/internal/health"
)

// NewServer returns an echo server with default configuration
func newServer() *echo.Echo {
	return echo.New()
}

// NewBackendsServer instantiates new application Echo server.
func NewBackendsServer(config *config.Config) (*echo.Echo, error) {
	e := newServer()

	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	 AllowOrigins: []string{"https://www.morphysm.com"},
	//	 AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

	const gitHost = "https://api.github.com"

	// TODO move to config
	const appID = 160183
	appClient, err := apps.NewClient(gitHost, config.Github.Key, appID)
	if err != nil {
		return nil, err
	}

	const owner = "morphysm"
	const installationID = 21534367
	var repoIDs = []int64{434540357, 440546811}
	token, err := appClient.GetAccessTokens(
		context.Background(),
		installationID,
		repoIDs)
	if err != nil {
		return nil, err
	}

	installationClient, err := installation.NewClient(gitHost, token, owner, installationID)
	if err != nil {
		return nil, err
	}

	githubHandler := glib.NewHandler(appClient, installationClient, installationID, config.Github.KudoLabel)

	// Logger
	e.Use(middleware.Logger())

	// GitHubRoutes endpoints exposed for Github requests.
	githubGroup := e.Group("/github")
	{
		GitHubRoutes(
			githubGroup, githubHandler,
		)
	}

	// Health endpoints exposed for heartbeat.
	healthGroup := e.Group("/health")
	{
		HealthRoutes(
			healthGroup, health.NewHandler(),
		)
	}

	return e, nil
}
