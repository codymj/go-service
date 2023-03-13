package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Respond converts a go value to json and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	// set status code for the request middleware
	_ = SetStatusCode(ctx, statusCode)

	// 204
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// marshal data
	jsn, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// send response
	if _, err = w.Write(jsn); err != nil {
		return err
	}

	return nil
}
