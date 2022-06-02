package server

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

// Server struct.
type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

// NewServer instantiates and setup new Server with new Echo server.
func NewServer(cfg *config.Config) (*Server, error) {
	nrApp, err := configureNewRelic(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't configure relic: %w", err)
	}

	echoServer := echo.New()

	echoServer.HideBanner = true
	echoServer.StdLogger.Printf(assets.Banner)

	// Middleware
	echoServer.Use(
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
		return nil, fmt.Errorf("can't create app client: %w", err)
	}

	// Get installations
	installations, err := appClient.GetInstallations(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't get installations: %w", err)
	}

	// Transform all installations to owner installationID map
	transformedInstallations := make(map[string]int64)
	for _, installation := range installations {
		transformedInstallations[installation.Account.Login] = installation.ID
	}

	// Create a new github client to fetch repo data
	installationClient, err := providers.NewInstallationClient(cfg.Github.Host, appClient, transformedInstallations, cfg.Github.WebhookSecret, cfg.Famed.Labels[config.FamedLabelKey].Name, cfg.RedTeamLogins)
	if err != nil {
		return nil, fmt.Errorf("can't create new github client: %w", err)
	}

	// Create a new GitHub handler handling gateway calls to GitHub
	githubHandler := github.NewHandler(installationClient)

	// Create the famed handler handling the famed business logic
	famedConfig := model.NewFamedConfig(cfg.Famed.Currency, cfg.Famed.Rewards, cfg.Famed.Labels, cfg.Famed.DaysToFix, cfg.Github.BotLogin)
	famedHandler := famed.NewHandler(appClient, installationClient, famedConfig, time.Now)

	// Start comment update interval
	ticker.NewTicker(time.Duration(cfg.Famed.UpdateFrequency)*time.Second, famedHandler.CleanState)

	// FamedRoutes endpoints exposed for Famed frontend client requests
	famedGroup := echoServer.Group("/famed")
	{
		FamedRoutes(
			famedGroup, famedHandler,
		)
	}

	// FamedAdminRoutes endpoints exposed for Famed admin requests
	famedAdminGroup := echoServer.Group("/admin", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
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
	healthGroup := echoServer.Group("/health")
	{
		HealthRoutes(
			healthGroup, health.NewHandler(),
		)
	}

	return &Server{
		echo: echoServer,
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

// Start starts the server and opens a new go routine that waits for the server to graceful shutdown it.
func (s *Server) Start() error {
	idleConnsClosed := make(chan struct{})

	// Waits for an interrupt to shutdown the server
	go func() {
		// The expected signals to turn off the server (CTRL+C/syscall interrupt)
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		defer stop()
		<-ctx.Done()

		// Does not accept any more requests, processes the remaining requests and stops the server
		log.Println("Requested shutdown in progress.. Press Ctrl+C again to force.")

		// Give 10 second to server to shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.echo.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		close(idleConnsClosed)
	}()

	// Start the server, the main thread will be blocked here
	if err := s.echo.Start(":" + s.cfg.App.Port); err != nil {
		close(idleConnsClosed)

		return fmt.Errorf("http server can't listen and serve: %w", err)
	}

	<-idleConnsClosed

	return nil
}
