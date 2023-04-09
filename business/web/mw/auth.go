package mw

import (
	"context"
	"errors"
	"fmt"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/codymj/go-service/business/sys/validate"
	"github.com/codymj/go-service/foundation/web"
	"net/http"
	"strings"
)

// Authenticate validates a JWT from the Authorization header.
func Authenticate(a *auth.Auth) web.Middleware {
	// function to be executed
	m := func(handler web.Handler) web.Handler {
		// create the handler that will be connected to mw chain
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// expecting: bearer <token>
			authStr := r.Header.Get("Authorization")

			// parse header
			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: bearer <token>")
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// validate the token is signed
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// add claims to the context for later retrieval
			ctx = auth.SetClaims(ctx, claims)

			// call next handler
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Authorize validates the authenticated user has at least one role.
func Authorize(roles ...string) web.Middleware {
	// function to be executed
	m := func(handler web.Handler) web.Handler {
		// create the handler that will be connected to mw chain
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// if context is missing, return failure
			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return validate.NewRequestError(
					fmt.Errorf("you are unauthorized for this action: no claims"),
					http.StatusForbidden)
			}

			// check if user is authorized for this action
			if !claims.Authorized(roles...) {
				return validate.NewRequestError(
					fmt.Errorf("you are unauthorized for this action: claims[%v] roles[%v]", claims.Roles, roles),
					http.StatusForbidden)
			}

			// call next handler
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
