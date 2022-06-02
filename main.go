package main

import (
	"log"

	"github.com/awnumar/memguard"

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

	// Instantiate and start the server
	if backendServer, err := server.NewServer(cfg); err != nil {
		log.Panic(err)
	} else {
		if err := backendServer.Start(); err != nil {
			log.Panic(err)
		}
	}
}
