package main

import (
	"fmt"
	"time"

	"github.com/awnumar/memguard"
	"github.com/phuslu/log"

	"github.com/morphysm/famed-github-backend/assets"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/server"
)

func main() {
	// Print the assets/banner.txt
	fmt.Println(assets.Banner)

	// Logger configuration
	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.Stamp,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: false,
		},
	}

	// Setup memguard https://pkg.go.dev/github.com/awnumar/memguard
	memguard.CatchInterrupt()
	defer memguard.Purge()

	// Load config
	cfg, err := config.NewConfig("config.json")
	if err != nil {
		log.Panic().Err(err).Msg("failed to load configuration")
	}

	// Instantiate the server
	backendServer, err := server.NewServer(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("failed to instantiate server")
	}

	if err := backendServer.Start(); err != nil {
		log.Panic().Err(err).Msg("failed to start server")
	}
}
