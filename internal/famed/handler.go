package famed

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/famed/model"
	"github.com/morphysm/famed-github-backend/internal/respositories/github/providers"
	"github.com/morphysm/famed-github-backend/pkg/sync"
)

type HTTPHandler interface {
	GetInstallations(c echo.Context) error
	GetTrackedIssues(c echo.Context) error

	GetBlueTeam(c echo.Context) error
	GetRedTeam(c echo.Context) error

	PostEvent(c echo.Context) error

	GetUpdateComments(c echo.Context) error

	CleanState()
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubAppClient          providers.AppClient
	githubInstallationClient providers.InstallationClient
	famedConfig              model.Config
	// now returns the current time
	// the time.Now function is not directly called to allow for testing
	now func() time.Time

	// TODO: investigate if this should be replace with a queue with a queue worker to avoid multiple blocked goroutines.
	issuesEventWG *sync.WaitGroups
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubAppClient providers.AppClient, githubInstallationClient providers.InstallationClient, famedConfig model.Config, now func() time.Time) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
		famedConfig:              famedConfig,
		now:                      now,
		issuesEventWG:            sync.NewWaitGroups(),
	}
}
