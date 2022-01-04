package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/config"
	"github.com/morphysm/kudos-github-backend/internal/server"
)

const (
	// http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=KudosBackend
	banner = `
   __ __        __            ______ __  __        __       ____             _        
  / //_/_ _____/ /__  _______/ ___(_) /_/ /  __ __/ /  ____/ __/__ _____  __(_)______ 
 / ,< / // / _  / _ \(_-<___/ (_ / / __/ _ \/ // / _ \/___/\ \/ -_) __/ |/ / / __/ -_)
/_/|_|\_,_/\_,_/\___/___/   \___/_/\__/_//_/\_,_/_.__/   /___/\__/_/  |___/_/\__/\__/ 

Go Backend
`
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare Main
	echoServer := prepareServer(config)

	// Start server.
	start(echoServer)
}

// Main setup
func prepareServer(config *config.Config) *echo.Echo {
	echoServer, echoServerErr := server.NewBackendsServer(config)
	if echoServerErr != nil {
		log.Fatal(echoServerErr)
	}

	return echoServer
}

// start an echo server with gracefully shutdown.
func start(e *echo.Echo) {
	// Start server for morphysm-service.
	go func() {
		e.HideBanner = true
		e.StdLogger.Printf(banner)
		if err := e.Start(":8080"); err != nil {
			log.Fatalf("shutting down the server. %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
