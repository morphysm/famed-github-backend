package server

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/github"
	"github.com/morphysm/kudos-github-backend/internal/health"
)

// GitHubRoutes defines endpoints exposed to serve relay calls to the Github api.
func GitHubRoutes(g *echo.Group, handler github.HTTPHandler) {
	g.GET("/installations", handler.GetInstallations)
	g.GET("/repos/:repo_name/labels", handler.GetLabels)
	g.GET("/repos/:repo_name/events", handler.GetEvents)
	g.GET("/repos/:repo_name/contributors", handler.GetContributors)
	g.GET("/repos/:repo_name/issues", handler.GetIssues)
	g.GET("/repos", handler.GetRepos)

	g.POST("/repos/:repo_name/issues/:issue_number/comments", handler.PostComment)
	g.POST("/webhooks/event", handler.PostEvent)
}

// HealthRoutes defines endpoints exposed to serve uses cases of infrastructure and customer support.
func HealthRoutes(g *echo.Group, handler health.HTTPHandler) {
	g.GET("", handler.GetHealth)
}
