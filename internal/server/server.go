package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/morphysm/kudos-github-backend/internal/client/app"
	"github.com/morphysm/kudos-github-backend/internal/client/installation"
	"github.com/morphysm/kudos-github-backend/internal/config"
	"github.com/morphysm/kudos-github-backend/internal/github"
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
	// TODO move to config
	const appID = "160183"
	appClient, err := app.NewClient("https://api.github.com", config.Github.Key, appID)
	if err != nil {
		return nil, err
	}

	const owner = "morphysm"
	const installationID = 21534367
	installationClient, err := installation.NewClient("https://api.github.com", appClient, owner, installationID)
	if err != nil {
		return nil, err
	}

	githubHandler := github.NewHandler(appClient, installationClient, config.Github.KudoLabel)

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
