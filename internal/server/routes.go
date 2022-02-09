package server

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/health"
)

// GitHubRoutes defines endpoints exposed to serve relay calls to the GitHub api.
func GitHubRoutes(g *echo.Group, handler famed.HTTPHandler) {
	g.GET("/repos/:repo_name/contributors", handler.GetContributors)

	g.POST("/webhooks/event", handler.PostEvent)

	g.POST("/repos/:repo_name/update", handler.UpdateComments)
}

// HealthRoutes defines endpoints exposed to serve uses cases of infrastructure and customer support.
func HealthRoutes(g *echo.Group, handler health.HTTPHandler) {
	g.GET("", handler.GetHealth)
}
