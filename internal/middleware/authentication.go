package middleware

import (
	"net/http"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
)

// AuthenticationMiddleware is the HTTP middleware to handle authentication.
type AuthenticationMiddleware struct {
	config *config.Config
}

// NewAuthenticationMiddleware returns a new AuthenticationMiddleware.
func NewAuthenticationMiddleware(config *config.Config) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		config: config,
	}
}

// Handle handles the HTTP request authentication.
func (m *AuthenticationMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			logger := log.CtxLogger(request.Context())

			if m.config.GetBool("config.authentication.enabled") {
				bearer := strings.TrimPrefix(request.Header.Get("authorization"), "Bearer ")
				secret := m.config.GetString("config.authentication.secret")
				if bearer != secret {
					logger.Warn().Msg("authentication failed")
					return echo.NewHTTPError(http.StatusUnauthorized, "authentication failure")
				}

				logger.Info().Msg("authentication success")
			}

			return next(c)
		}
	}
}
