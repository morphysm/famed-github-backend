// Package devtoolkit gathers several useful dependencies in different places in the code and avoids making global variables. This is a dependency injection.
package devtoolkit

import (
	"github.com/morphysm/famed-github-backend/internal/config"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit/buildinfo"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit/userdirs"

	"github.com/rotisserie/eris"

	"github.com/phuslu/log"
)

// DevToolkit store the various dependencies that will be used during injection.
type DevToolkit struct {
	// Logger represents the logging system.
	Logger *log.Logger
	// Config represents all configuration.
	Config *config.Config
	// SentryClient allows to send errors to Sentry.
	SentryClient *sentry.Client
	// UserDirs gives the paths of the folders respecting the XDG standards.
	UserDirs *userdirs.UserDirs
	// BuildInfo brings together all information about the current build.
	BuildInfo *buildinfo.BuildInfo
}

// NewDevToolkit instantiates a new DevToolkit for the application. This function should only be called once.
func NewDevToolkit() (toolkit *DevToolkit, err error) {
	toolkit = &DevToolkit{}

	// Logging system initialization
	toolkit.Logger = &log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.Stamp,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: false,
		},
	}

	// Generation of build information
	toolkit.BuildInfo, err = buildinfo.NewBuildInfo()
	if err != nil {
		return nil, eris.Wrap(err, "failed to instantiate build information")
	}

	// New configuration, based on environment variables and file
	toolkit.Config, err = config.Load()
	if err != nil {
		return nil, eris.Wrap(err, "failed to load configuration")
	}

	// Creating a new sentry client and checking the DSN
	// TODO: CHANGE THE DSN!
	toolkit.SentryClient, err = sentry.NewClient(sentry.ClientOptions{
		Dsn:              "https://foo@bar.ingest.sentry.io/TEST",
		AttachStacktrace: true,
		Release:          buildinfo.ProgramName + "@" + toolkit.BuildInfo.Version.String(),
		Environment:      toolkit.BuildInfo.Version.Prerelease(),
	})
	if err != nil {
		return nil, eris.Wrap(err, "failed to instantiate Sentry client")
	}

	toolkit.UserDirs, err = userdirs.NewUserDirs(buildinfo.ProgramName)
	if err != nil {
		return nil, eris.Wrap(err, "failed to instantiate user directories")
	}

	return toolkit, nil
}
