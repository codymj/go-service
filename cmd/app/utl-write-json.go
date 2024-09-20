package main

import (
	"encoding/json"
	"net/http"
)

// envelope for wrapping responses.
type envelope map[string]any

// writeJson writes JSON responses to the client.
func (a *app) writeJson(
	w http.ResponseWriter,
	headers http.Header,
	status int,
	data envelope,
) error {
	// Encode data to JSON.
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body = append(body, '\n')

	// Append headers.
	for key, val := range headers {
		w.Header()[key] = val
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Send response
	_, err = w.Write(body)
	if err != nil {
		return err
	}

	return nil
}
