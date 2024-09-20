package main

import (
	"net/http"
)

// Handler for GET /health
func (a *app) getHealth(w http.ResponseWriter, _ *http.Request) {
	// Construct body.
	data := map[string]string{
		"status":      "ok",
		"environment": a.config.env,
		"version":     version,
	}

	// Send response.
	err := a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
