package main

import (
	"fmt"
	"go-service.codymj.io/internal/data"
	"net/http"
	"time"
)

// getUsersId is a handler for GET /users/:id.
func (a *app) getUsersId(w http.ResponseWriter, r *http.Request) {
	// Parse id path parameter.
	id, err := a.parseIdParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	location := "Tampa, FL"
	user := data.User{
		Id:          id,
		UserName:    "codymj",
		FirstName:   "Cody",
		LastName:    "Johnson",
		Email:       "codyj@protonmail.com",
		Password:    "password123",
		Location:    &location,
		DateOfBirth: time.Date(1987, 3, 8, 0, 0, 0, 0, time.UTC),
		Created:     time.Now().UTC(),
		LastSeen:    time.Now().UTC(),
	}

	// Call service.
	err = a.writeJson(w, nil, http.StatusOK, envelope{"user": user})
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// postUsers is a handler for POST /users.
func (a *app) postUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate request.

	// Call service.
	_, err := fmt.Fprintln(w, "Create a new user.")
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
