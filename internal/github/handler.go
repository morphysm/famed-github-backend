package github

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/client"
)

type HTTPHandler interface {
	GetInstallations(c echo.Context) error
	GetAccessTokens(c echo.Context) error
}

// githubHandler represents the handler for the github endpoints.
type githubHandler struct {
	githubClient client.Client
}

// NewHandler returns a pointer to the github handler.
func NewHandler(githubClient client.Client) HTTPHandler {
	return &githubHandler{githubClient: githubClient}
}
