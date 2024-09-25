package handler

import (
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

// Dashboard is the HTTP handler to render the dashboard.
type DashboardHandler struct {
	config *config.Config
}

// NewDashboardHandler returns a new DashboardHandler.
func NewDashboardHandler(config *config.Config) *DashboardHandler {
	return &DashboardHandler{
		config: config,
	}
}

// Handle handles the HTTP request.
func (h *DashboardHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "dashboard.html", map[string]any{
			"title": h.config.GetString("config.dashboard.title"),
		})
	}
}
