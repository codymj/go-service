package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const v1 = "/v1"

func (a *app) routes() http.Handler {
	// Initialize router instance.
	router := httprouter.New()

	// Application routes.
	router.HandlerFunc(http.MethodGet, v1+"/health", a.getHealth)

	// User routes.
	//router.HandlerFunc(http.MethodGet, v1+"/users", a.getUsers)
	router.HandlerFunc(http.MethodGet, v1+"/users/:id", a.getUsersId)
	router.HandlerFunc(http.MethodPost, v1+"/users", a.postUsers)

	return router
}
