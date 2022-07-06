package server

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/github"
	"github.com/morphysm/famed-github-backend/internal/health"
)

// FamedRoutes defines endpoints exposed to serve famed api endpoints.
func FamedRoutes(g *echo.Group, handler famed.HTTPHandler) {
	g.GET("/repos/:owner/:repo_name/contributors", handler.GetBlueTeam)
	g.GET("/repos/:owner/:repo_name/redteam", handler.GetRedTeam)

	g.POST("/webhooks/event", handler.PostEvent)

	g.POST("/repos/:owner/:repo_name/update", handler.GetUpdateComments)

	g.GET("/:owner/eligableissues", handler.GetElligableIssues)
	g.GET("/:contributor/rewards", handler.GetContributorRewards)
}

func FamedAdminRoutes(g *echo.Group, famedHandler famed.HTTPHandler, githubHandler github.HTTPHandler) {
	g.GET("/installations", famedHandler.GetInstallations)
	g.GET("/trackedissues", famedHandler.GetTrackedIssues)
	g.GET("/ratelimit/:owner", githubHandler.GetRateLimit)
}

// HealthRoutes defines endpoints exposed to serve uses cases of infrastructure and customer support.
func HealthRoutes(g *echo.Group, handler health.HTTPHandler) {
	g.GET("", handler.GetHealth)
}
