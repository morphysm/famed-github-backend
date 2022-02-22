package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/morphysm/famed-github-backend/internal/client/apps"
	"github.com/morphysm/famed-github-backend/internal/client/currency"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/health"
)

// NewServer returns an echo server with default configuration
func newServer() *echo.Echo {
	return echo.New()
}

// NewBackendsServer instantiates new application Echo server.
func NewBackendsServer(config *config.Config) (*echo.Echo, error) {
	e := newServer()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://www.famed.morphysm.com", "https://famed.morphysm.com"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	currencyClient := currency.NewCurrencyClient(config.Currency.Host)

	appClient, err := apps.NewClient(config.Github.Host, config.Github.Key, config.Github.AppID)
	if err != nil {
		return nil, err
	}

	installationClient, err := installation.NewClient(config.Github.Host, appClient, config.Github.InstallationID, config.Github.RepoIDs, config.Github.Owner)
	if err != nil {
		return nil, err
	}

	famedConfig := famed.Config{
		Label:     config.Famed.Label,
		Currency:  config.Famed.Currency,
		Rewards:   config.Famed.Rewards,
		BotUserID: config.Github.BotID,
	}
	githubHandler := famed.NewHandler(installationClient, currencyClient, &config.Github.WebhookSecret, config.Github.InstallationID, famedConfig)

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
