package main

import (
	"fmt"
	"net/http"
)

// logError logs error message.
func (a *app) logError(r *http.Request, err error) {
	a.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

// sendError sends a JSON-formatted error response to client.
func (a *app) sendError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	msg any,
) {
	err := a.writeJson(w, nil, status, envelope{"error": msg})
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// send400 sends a JSON-formatted 400 response to client.
func (a *app) send400(w http.ResponseWriter, r *http.Request, err error) {
	a.sendError(w, r, http.StatusBadRequest, err.Error())
}

// send404 sends a JSON-formatted 404 response to client.
func (a *app) send404(w http.ResponseWriter, r *http.Request) {
	msg := "The requested resource could not be found."
	a.sendError(w, r, http.StatusNotFound, msg)
}

// send405 sends a JSON-formatted 405 response to client.
func (a *app) send405(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf(
		"The %s method is not supported for this resource.",
		r.Method,
	)
	a.sendError(w, r, http.StatusMethodNotAllowed, msg)
}

// send500 sends a JSON-formatted 500 response to client.
func (a *app) send500(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	msg := "An internal error occurred and could not process your request."
	a.sendError(w, r, http.StatusInternalServerError, msg)
}
