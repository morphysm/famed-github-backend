package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/awnumar/memguard"
	"github.com/morphysm/famed-github-backend/assets"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit/buildinfo"
	"github.com/morphysm/famed-github-backend/internal/subcommand"
	"github.com/phuslu/log"
	"time"
)

// Arguments are all the possible subcommands, arguments and flags that can be sent to the application.
type Arguments struct {
	Server *Server `arg:"subcommand:server" help:"Start the server (default)"` // Server is the subcommand that starts the server.
}

// Server subcommand starts the server.
type Server struct{}

// Version prints build information (--version argument).
func (Arguments) Version() string {
	buildinfo, err := buildinfo.NewBuildInfo()
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return fmt.Sprintf("%+v", *buildinfo)
}

// Description returns the software description (--help argument).
func (Arguments) Description() string {
	return "☄ " + buildinfo.ProjectName + " is a LOREM foobar.️\n" +
		"🌐 " + buildinfo.ProjectWebsite
}

func main() {
	// Define default arguments
	arguments := &Arguments{
		Server: nil,
	}

	// Parse the arguments and set server as default sub-command.
	if arg.MustParse(arguments).Subcommand() == nil {
		arguments.Server = &Server{}
	}

	// Setup essential components (log, config and sentry)
	devtoolkit, err := devtoolkit.NewDevToolkit()
	if err != nil {
		log.Panic().Err(err).Msg("failed to setup essential components")
	}

	// Print the assets/banner.txt
	fmt.Println(assets.Banner)

	// Set logger level
	devtoolkit.Logger.Level = log.ErrorLevel

	// Setup memguard https://pkg.go.dev/github.com/awnumar/memguard
	memguard.CatchInterrupt()
	defer memguard.Purge()

	// Check and run server subcommand
	if arguments.Server != nil {
		if serverSubCmd, err := subcommand.NewServer(devtoolkit); err != nil {
			devtoolkit.Logger.Error().Err(err).Msg("can't initialize server subcommand")
		} else {
			if err := serverSubCmd.Start(); err != nil {
				devtoolkit.Logger.Error().Err(err).Msg("can't start server subcommand")
			}
		}
	}

	// Close sentry
	devtoolkit.SentryClient.Flush(time.Second + time.Second)
}
