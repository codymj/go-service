package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// parseIdParameter parses an ID parameter from the requested path.
func (a *app) parseIdParameter(r *http.Request) (int64, error) {
	parameters := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(parameters.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID path parameter")
	}

	return id, nil
}

// writeJson writes JSON responses to the client.
func (a *app) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
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
