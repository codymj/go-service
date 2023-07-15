package usergroup

import (
	"context"
	"errors"
	"fmt"
	"github.com/codymj/go-service/business/core/user"
	userstore "github.com/codymj/go-service/business/data/store/user"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/codymj/go-service/business/sys/validate"
	"github.com/codymj/go-service/foundation/web"
	"net/http"
	"strconv"
)

// Handlers managers the set of user endpoints.
type Handlers struct {
	Auth *auth.Auth
	User user.Core
}

// Query returns a list of users with pagination.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get page path parameter
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("usergroup.Query(): %s", page), http.StatusBadRequest)
	}

	// get rows path parameter
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return validate.NewRequestError(fmt.Errorf("usergroup.Query(): %s", rows), http.StatusBadRequest)
	}

	// get users
	users, err := h.User.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("usergroup.Query(): %w", err)
	}

	return web.Respond(ctx, w, users, http.StatusOK)
}

// QueryById returns a user by its ID.
func (h Handlers) QueryById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get claims from context
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("claims missing from context")
	}

	// get user id path param
	id := web.Param(r, "id")

	// get user
	usr, err := h.User.QueryById(ctx, claims, id)
	if err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidId:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("usergroup.QueryById(): %w", err)
		}
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

// Create adds a new user to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get context values
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	// construct a new user
	var nu userstore.NewUser
	if err = web.Decode(r, &nu); err != nil {
		return fmt.Errorf("usergroup.Create(): %w", err)
	}

	// create user
	usr, err := h.User.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("usergroup.Create(): %w", err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Update updates a user in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get context values
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	// get claims
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("claims missing from context")
	}

	// construct updated user
	var uu userstore.UpdateUser
	if err = web.Decode(r, &uu); err != nil {
		return fmt.Errorf("usergroup.Update(): %w", err)
	}

	// get user id from path param
	id := web.Param(r, "id")

	// update user
	if err = h.User.Update(ctx, claims, id, uu, v.Now); err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidId:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("usergroup.Update(): %w", err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a user from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get claims
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("claims missing from context")
	}

	// get user id from path param
	id := web.Param(r, "id")

	// delete user
	if err = h.User.Delete(ctx, claims, id); err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidId:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("usergroup.Delete(): %w", err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Token provides an API token for the authenticated user.
func (h Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get context values
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	// get email, password from request
	email, pass, ok := r.BasicAuth()
	if !ok {
		err = errors.New("must provide email and password in Basic auth")
		return validate.NewRequestError(err, http.StatusUnauthorized)
	}

	// authenticate user
	claims, err := h.User.Authenticate(ctx, v.Now, email, pass)
	if err != nil {
		switch validate.Cause(err) {
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrAuthFailure:
			return validate.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("usergroup.Token(): %w", err)
		}
	}

	// generate token
	var token struct {
		Token string `json:"token"`
	}
	token.Token, err = h.Auth.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("usergroup.Token(): %w", err)
	}

	return web.Respond(ctx, w, token, http.StatusOK)
}
