package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/morphysm/famed-github-backend/internal/client/app"
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

	// Create new app client to fetch installations and installation tokens.
	appClient, err := app.NewClient(config.Github.Host, config.Github.Key, config.Github.AppID)
	if err != nil {
		return nil, err
	}

	// Get installations
	installations, err := appClient.GetInstallations(context.Background())
	if err != nil {
		return nil, err
	}

	// Transform all installations to owner installationID map
	transformedInstallations := make(map[string]int64)
	for _, installation := range installations {
		transformedInstallations[*installation.Account.Login] = *installation.ID
	}

	// Create new installation client to fetch repo data
	installationClient, err := installation.NewClient(config.Github.Host, appClient, transformedInstallations)
	if err != nil {
		return nil, err
	}

	// Create
	famedConfig := famed.Config{
		Currency:  config.Famed.Currency,
		Rewards:   config.Famed.Rewards,
		Labels:    config.Famed.Labels,
		BotUserID: config.Github.BotID,
	}
	famedHandler := famed.NewHandler(installationClient, currencyClient, &config.Github.WebhookSecret, famedConfig)

	// Logger
	e.Use(middleware.Logger())

	// GitHubRoutes endpoints exposed for Github requests.
	famedGroup := e.Group("/famed")
	{
		FamedRoutes(
			famedGroup, famedHandler,
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
