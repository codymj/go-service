package main

import (
	"fmt"
	"go-service.codymj.io/internal/data"
	"net/http"
	"time"
)

// getUsersId is a handler for GET /v1/users/:id.
func (a *app) getUsersId(w http.ResponseWriter, r *http.Request) {
	// Parse id path parameter.
	id, err := a.parseIdParam(r)
	if err != nil {
		a.send404(w, r)
		return
	}

	location := "Tampa, FL"
	user := data.User{
		Id:          id,
		Username:    "codymj",
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
		a.send500(w, r, err)
	}
}

// postUsers is a handler for POST /v1/users.
func (a *app) postUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate request body with a JSON validator library.

	// Struct for decoding request.
	var body struct {
		Username    string    `json:"username"`
		FirstName   string    `json:"firstName"`
		LastName    string    `json:"lastName"`
		Email       string    `json:"email"`
		Password    string    `json:"password"`
		Location    string    `json:"location,omitempty"`
		DateOfBirth time.Time `json:"dateOfBirth"`
	}

	// Decode request body.
	err := a.readJson(w, r, &body)
	if err != nil {
		a.send400(w, r, err)
		return
	}

	// Call service.
	_, err = fmt.Fprintf(w, "%+v\n%", body)
	if err != nil {
		a.send500(w, r, err)
	}
}
