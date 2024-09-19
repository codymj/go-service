package main

import (
	"fmt"
	"net/http"
)

// getUsersId is a handler for GET /users/:id.
func (a *app) getUsersId(w http.ResponseWriter, r *http.Request) {
	// Parse id path parameter.
	id, err := a.parseIdParameter(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Call service.
	_, err = fmt.Fprintf(w, "Get user %d\n", id)
	if err != nil {
		a.logger.Error(err.Error())
	}
}

// postUsers is a handler for POST /users.
func (a *app) postUsers(w http.ResponseWriter, _ *http.Request) {
	// TODO: Validate request.

	// Call service.
	_, err := fmt.Fprintln(w, "Create a new user.")
	if err != nil {
		a.logger.Error(err.Error())
	}
}
