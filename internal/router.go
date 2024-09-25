package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"go-service.codymj.io/internal/handler"
	"go-service.codymj.io/internal/handler/user"
	"go-service.codymj.io/internal/middleware"
	"go.uber.org/fx"
)

// Router is used to register the application HTTP middlewares and handlers.
func Router() fx.Option {
	return fx.Options(
		// Authentication middleware:
		fxhttpserver.AsMiddleware(middleware.NewAuthenticationMiddleware, fxhttpserver.GlobalUse),

		// Dashboard handler:
		fxhttpserver.AsHandler("GET", "", handler.NewDashboardHandler),

		// Users CRUD handlers:
		fxhttpserver.AsHandlersGroup(
			"/users",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", user.NewListUsersHandler),
				// fxhttpserver.NewHandlerRegistration("GET", "/:id", user.NewGetUserHandler),
				// fxhttpserver.NewHandlerRegistration("POST", "", user.NewCreateUserHandler),
				// fxhttpserver.NewHandlerRegistration("DELETE", "/:id", user.NewDeleteUserHandler),
			},
		),
	)
}
