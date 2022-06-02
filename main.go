package main

import (
	"context"
	"github.com/morphysm/famed-github-backend/assets"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/awnumar/memguard"
	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/server"
)

func main() {
	// Setup memguard https://pkg.go.dev/github.com/awnumar/memguard
	memguard.CatchInterrupt()
	defer memguard.Purge()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Panic(err)
	}

	// Prepare server
	echoServer := prepareServer(cfg)

	// Start server
	start(echoServer, cfg)
}

// Main setup
func prepareServer(cfg *config.Config) *echo.Echo {
	echoServer, echoServerErr := server.NewBackendServer(cfg)
	if echoServerErr != nil {
		log.Fatal(echoServerErr)
	}

	return echoServer
}

// start an echo server with gracefully shutdown.
func start(e *echo.Echo, cfg *config.Config) {
	// Start server for famed backend.
	go func() {
		e.HideBanner = true
		e.StdLogger.Printf(assets.Banner)

		if err := e.Start(":" + cfg.App.Port); err != nil {
			log.Fatalf("shutting down the server. %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Panic(err)
	}
}
