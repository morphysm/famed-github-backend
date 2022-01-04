package server

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/github"
	"github.com/morphysm/kudos-github-backend/internal/health"
)

// GitHubRoutes defines endpoints exposed to serve relay calls to the Github api.
func GitHubRoutes(g *echo.Group, handler github.HTTPHandler) {
	g.GET("/installations", handler.GetInstallations)
	g.GET("/repos/:repo_id/labels", handler.GetLabels)
	g.GET("/repos/", handler.GetRepos)
}

// HealthRoutes defines endpoints exposed to serve uses cases of infrastructure and customer support.
func HealthRoutes(g *echo.Group, handler health.HTTPHandler) {
	g.GET("", handler.GetHealth)
}
