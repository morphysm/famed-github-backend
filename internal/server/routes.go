package server

import (
	"github.com/labstack/echo/v4"

	"github.com/morphysm/kudos-github-backend/internal/health"
)


// HealthRoutes defines endpoints exposed to serve uses cases of infrastructure and customer support.
func HealthRoutes(g *echo.Group, handler health.HTTPHandler) {
	g.GET("", handler.Health)
}
