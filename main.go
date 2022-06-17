package main

import (
	"fmt"
	"github.com/morphysm/famed-github-backend/internal/otherconfig"

	"github.com/awnumar/memguard"
	"github.com/morphysm/famed-github-backend/assets"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/server"
	"github.com/phuslu/log"
)

func main() {

	config, err := otherconfig.NewConfig("./config.yaml")
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		return
	}

	fmt.Println(config.App.Host)

	return

	// Print the assets/banner.txt
	fmt.Println(assets.Banner)

	// Setup essential components (log, config and sentry)
	devtoolkit, err := devtoolkit.NewDevToolkit()
	if err != nil {
		log.Panic().Err(err).Msg("failed to setup essential components")
	}

	// Setup memguard https://pkg.go.dev/github.com/awnumar/memguard
	memguard.CatchInterrupt()
	defer memguard.Purge()

	// Instantiate the server
	backendServer, err := server.NewServer(devtoolkit)
	if err != nil {
		log.Panic().Err(err).Msg("failed to instantiate server")
	}

	if err := backendServer.Start(); err != nil {
		log.Panic().Err(err).Msg("failed to start server")
	}
}
