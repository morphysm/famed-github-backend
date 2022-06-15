package main

import (
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/morphysm/famed-github-backend/assets"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/server"
	"github.com/phuslu/log"
)

func main() {
	// Print the assets/banner.txt
	fmt.Println(assets.Banner)

	// Instantiate essential components (log, config, etc.)
	devtoolkit, err := devtoolkit.NewDevToolkit()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to instantiate essential components")
	}

	// Setup memguard https://pkg.go.dev/github.com/awnumar/memguard
	memguard.CatchInterrupt()
	defer memguard.Purge()

	// Instantiate the server
	backendServer, err := server.NewServer(devtoolkit.Config)
	if err != nil {
		log.Panic().Err(err).Msg("failed to instantiate server")
	}

	if err := backendServer.Start(); err != nil {
		log.Panic().Err(err).Msg("failed to start server")
	}
}
