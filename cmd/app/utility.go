package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// parseIdParam parses an ID parameter from the requested path.
func (a *app) parseIdParam(r *http.Request) (int64, error) {
	parameters := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(parameters.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID path parameter")
	}

	return id, nil
}

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

// readJson reads JSON request bodies from the client.
func (a *app) readJson(w http.ResponseWriter, r *http.Request, dest any) error {
	// Set maximum bytes allowed in request body.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize decoder.
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the request body.
	err := decoder.Decode(dest)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError
		var maxBytesErr *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf(
				"body contains malformed JSON at character %d",
				syntaxErr.Offset,
			)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains malformed JSON")
		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeErr.Field,
				)
			}
			return fmt.Errorf(
				"body contains incorrect JSON type at character %d",
				unmarshalTypeErr.Offset,
			)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesErr):
			return fmt.Errorf(
				"body must not be larged than %d bytes",
				maxBytesErr.Limit,
			)
		case errors.As(err, &invalidUnmarshalErr):
			panic(err)
		default:
			return err
		}
	}

	// Call decode again.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON structure")
	}

	return nil
}
