package famed

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/client/github"
	"github.com/morphysm/famed-github-backend/pkg/sync"
)

type HTTPHandler interface {
	GetInstallations(c echo.Context) error
	GetTrackedIssues(c echo.Context) error
	GetRateLimit(c echo.Context) error

	GetContributors(c echo.Context) error
	GetRedTeam(c echo.Context) error

	PostEvent(c echo.Context) error

	GetUpdateComments(c echo.Context) error

	CleanState()
}

// githubHandler represents the handler for the GitHub endpoints.
type githubHandler struct {
	githubAppClient          github.AppClient
	githubInstallationClient github.InstallationClient

	// TODO: investigate if this should be replace with a queue with a queue worker to avoid multiple blocked goroutines.
	issuesEventWG *sync.WaitGroups

	famedConfig Config
}

// NewHandler returns a pointer to the GitHub handler.
func NewHandler(githubAppClient github.AppClient, githubInstallationClient github.InstallationClient, famedConfig Config) HTTPHandler {
	return &githubHandler{
		githubAppClient:          githubAppClient,
		githubInstallationClient: githubInstallationClient,
		issuesEventWG:            sync.NewWaitGroups(),
		famedConfig:              famedConfig,
	}
}
