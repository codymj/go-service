package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go-service.codymj.io/internal/repository"
	"go-service.codymj.io/internal/service"
)

// ListUsersHandler is the HTTP handler to list all users.
type ListUsersHandler struct {
	service *service.UserService
}

// NewListUsersHandler returns a new ListUsersHandler.
func NewListUsersHandler(service *service.UserService) *ListUsersHandler {
	return &ListUsersHandler{
		service: service,
	}
}

// Handle handles the HTTP request.
func (h *ListUsersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse query parameters from request.
		var username sql.NullString
		username.String = c.QueryParam("username")

		var email sql.NullString
		email.String = c.QueryParam("email")

		params := repository.UserRepositoryFindAllParams{
			Username: username,
			Email:    email,
		}

		// Call service to get users.
		users, err := h.service.List(c.Request().Context(), params)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("cannot list users: %v", err),
			)
		}

		return c.JSON(http.StatusOK, users)
	}
}
