package server

import (
	"context"
	"crypto/subtle"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rotisserie/eris"

	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/github"
	"github.com/morphysm/famed-github-backend/internal/health"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/providers"
	"github.com/morphysm/famed-github-backend/pkg/ticker"
)

// Server represents the HTTP server single instance.
type Server struct {
	echo       *echo.Echo
	devToolKit *devtoolkit.DevToolkit
}

// NewServer instantiates and sets up a new server using the echo web framework.
func NewServer(devToolKit *devtoolkit.DevToolkit) (*Server, error) {
	nrApp, err := configureNewRelic(devToolKit.Config)
	if err != nil {
		return nil, eris.Wrap(err, "failed to configure relic")
	}

	echoServer := echo.New()

	echoServer.HideBanner = true

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
	appClient, err := providers.NewAppClient(devToolKit.Config.Github.Host, devToolKit.Config.Github.AppID, devToolKit.Config.Github.KeyEnclave)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create app client")
	}

	// Get installations
	installations, err := appClient.GetInstallations(context.Background())
	if err != nil {
		return nil, eris.Wrap(err, "failed to get installations")
	}

	// Transform all installations to owner installationID map
	transformedInstallations := make(map[string]int64)
	for _, installation := range installations {
		transformedInstallations[installation.Account.Login] = installation.ID
	}

	// Create a new github client to fetch repo data
	installationClient, err := providers.NewInstallationClient(devToolKit.Config.Github.Host, appClient, transformedInstallations, devToolKit.Config.Github.WebhookSecret, devToolKit.Config.Famed.Labels[config.FamedLabelKey].Name, devToolKit.Config.RedTeamLogins)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create new github client")
	}

	// Create a new GitHub handler handling gateway calls to GitHub
	githubHandler := github.NewHandler(installationClient)

	// Create the famed handler handling the famed business logic
	famedConfig := model.NewFamedConfig(devToolKit.Config.Famed.Currency, devToolKit.Config.Famed.Rewards, devToolKit.Config.Famed.Labels, devToolKit.Config.Famed.DaysToFix, devToolKit.Config.Github.BotLogin)
	famedHandler := famed.NewHandler(appClient, installationClient, famedConfig, time.Now)

	// Start comment update interval
	ticker.NewTicker(time.Duration(devToolKit.Config.Famed.UpdateFrequency)*time.Second, famedHandler.CleanState)

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
		if subtle.ConstantTimeCompare([]byte(username), []byte(devToolKit.Config.Admin.Username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(devToolKit.Config.Admin.Password)) == 1 {
			return true, nil
		}

		return false, nil
	}))
	{
		FamedAdminRoutes(
			famedAdminGroup, famedHandler, githubHandler, devToolKit.Monitor,
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
		echo:       echoServer,
		devToolKit: devToolKit,
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

// Start starts a new go routine that allows to gracefully shut down the server
func (s *Server) Start() error {
	idleConnsClosed := make(chan struct{})

	// Waits for an interrupt to shutdown the server
	go func() {
		// The expected signals to turn off the server (CTRL+C/syscall interrupt)
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		defer stop()
		<-ctx.Done()

		// Does not accept any more requests, processes the remaining requests and stops the server
		s.devToolKit.Logger.Info().Msg("Requested shutdown in progress.. Press Ctrl+C again to force.")

		// Give 10 second to server to shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.echo.Shutdown(ctx); err != nil {
			s.devToolKit.Logger.Fatal().Err(err).Msg("failed to gracefully shutdown server")
		}

		close(idleConnsClosed)
	}()

	// Start the server, the main thread will be blocked here
	if err := s.echo.Start(net.JoinHostPort("", s.devToolKit.Config.App.Port)); err != nil {
		close(idleConnsClosed)

		return eris.Wrap(err, "failed to listen and serve http server")
	}

	<-idleConnsClosed

	return nil
}
