package main

import (
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
