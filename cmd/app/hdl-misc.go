package main

import (
	"net/http"
)

// getHealth handles GET /v1/health.
func (a *app) getHealth(w http.ResponseWriter, r *http.Request) {
	// Construct body.
	data := envelope{
		"status": "ok",
		"app": map[string]string{
			"environment": a.config.env,
			"version":     version,
		},
	}

	// Send response.
	err := a.writeJson(w, nil, http.StatusOK, data)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
