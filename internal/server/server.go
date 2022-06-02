package server

import (
	"context"
	"crypto/subtle"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/morphysm/famed-github-backend/assets"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/github"
	"github.com/morphysm/famed-github-backend/internal/health"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
	"github.com/morphysm/famed-github-backend/pkg/ticker"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

// NewBackendServer instantiates new application Echo server.
func NewServer(cfg *config.Config) (*Server, error) {
	nrApp, err := configureNewRelic(cfg)
	if err != nil {
		return nil, err
	}

	e := echo.New()

	e.HideBanner = true
	e.StdLogger.Printf(assets.Banner)

	// Middleware
	e.Use(
		nrecho.Middleware(nrApp),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://www.famed.morphysm.com", "https://famed.morphysm.com"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}),
		middleware.Logger(),
	)

	// Create new app client to fetch installations and github tokens.
	appClient, err := providers.NewAppClient(cfg.Github.Host, cfg.Github.AppID, cfg.Github.KeyEnclave)
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
		transformedInstallations[installation.Account.Login] = installation.ID
	}

	// Create a new github client to fetch repo data
	installationClient, err := providers.NewInstallationClient(cfg.Github.Host, appClient, transformedInstallations, cfg.Github.WebhookSecret, cfg.Famed.Labels[config.FamedLabelKey].Name, cfg.RedTeamLogins)
	if err != nil {
		return nil, err
	}

	// Create a new GitHub handler handling gateway calls to GitHub
	githubHandler := github.NewHandler(installationClient)

	// Create the famed handler handling the famed business logic
	famedConfig := model.NewFamedConfig(cfg.Famed.Currency, cfg.Famed.Rewards, cfg.Famed.Labels, cfg.Famed.DaysToFix, cfg.Github.BotLogin)
	famedHandler := famed.NewHandler(appClient, installationClient, famedConfig, time.Now)

	// Start comment update interval
	ticker.NewTicker(time.Duration(cfg.Famed.UpdateFrequency)*time.Second, famedHandler.CleanState)

	// FamedRoutes endpoints exposed for Famed frontend client requests
	famedGroup := e.Group("/famed")
	{
		FamedRoutes(
			famedGroup, famedHandler,
		)
	}

	// FamedAdminRoutes endpoints exposed for Famed admin requests
	famedAdminGroup := e.Group("/admin", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Use of constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Admin.Username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Admin.Password)) == 1 {
			return true, nil
		}
		return false, nil
	}))
	{
		FamedAdminRoutes(
			famedAdminGroup, famedHandler, githubHandler,
		)
	}

	// Health endpoints exposed for heartbeat
	healthGroup := e.Group("/health")
	{
		HealthRoutes(
			healthGroup, health.NewHandler(),
		)
	}

	return &Server{
		echo: e,
		cfg:  cfg,
	}, nil
}

func configureNewRelic(cfg *config.Config) (*newrelic.Application, error) {
	if !cfg.NewRelic.Enabled {
		return newrelic.NewApplication(
			newrelic.ConfigEnabled(cfg.NewRelic.Enabled),
		)
	}
	return newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.NewRelic.Name),
		newrelic.ConfigLicense(cfg.NewRelic.Key),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigEnabled(cfg.NewRelic.Enabled),
	)
}

func (s Server) Start() {
	// Start server for famed backend.
	go func() {
		if err := s.echo.Start(":" + s.cfg.App.Port); err != nil {
			log.Fatalf("shutting down the server. %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		log.Panic(err)
	}
}
