package server

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/health"
	"github.com/morphysm/kudos-github-backend/internal/kudo"
)

// GitHubRoutes defines endpoints exposed to serve relay calls to the GitHub api.
func GitHubRoutes(g *echo.Group, handler kudo.HTTPHandler) {
	g.GET("/repos/:repo_name/contributors", handler.GetContributors)

	g.POST("/webhooks/event", handler.PostEvent)
}

// HealthRoutes defines endpoints exposed to serve uses cases of infrastructure and customer support.
func HealthRoutes(g *echo.Group, handler health.HTTPHandler) {
	g.GET("", handler.GetHealth)
}
