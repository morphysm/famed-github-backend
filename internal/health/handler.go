package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPHandler interface {
	GetHealth(c echo.Context) error
}

// healthHandler represents the handler for the health endpoints.
type healthHandler struct{}

// NewHandler returns a pointer to the health handler.
func NewHandler() HTTPHandler {
	return &healthHandler{}
}

// healthStatus represents the data exposed on health endpoint.
type healthStatus struct {
	// Version of service
	Version string `json:"version"`
}

// GetHealth returns a json with the service version
func (healthHandler *healthHandler) GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, healthStatus{
		Version: "0.0.1",
	})
}
