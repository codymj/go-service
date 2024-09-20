package main

import (
	"fmt"
	"net/http"
)

// logError logs error message.
func (a *app) logError(r *http.Request, err error) {
	a.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

// errorResponse sends a JSON-formatted error response to client.
func (a *app) errorResponse(
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

// serverErrorResponse sends a JSON-formatted 500 response to client.
func (a *app) serverErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	a.logError(r, err)

	msg := "An internal error occurred and could not process your request."
	a.errorResponse(w, r, http.StatusInternalServerError, msg)
}

// notFoundResponse sends a JSON-formatted 404 response to client.
func (a *app) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "The requested resource could not be found."
	a.errorResponse(w, r, http.StatusNotFound, msg)
}

// notAllowedResponse sends a JSON-formatted 405 response to client.
func (a *app) notAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf(
		"The %s method is not supported for this resource.",
		r.Method,
	)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}
