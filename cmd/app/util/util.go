package util

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	ContentType = "Content-Type"
	JsonHeader  = "application/json"
)

func ParseQueryString(s string) map[string]string {
	params := make(map[string]string)

	splitQueries := strings.Split(s, "&")
	for _, query := range splitQueries {
		keyval := strings.Split(query, "=")
		key := keyval[0]
		val := keyval[1]
		params[key] = val
	}

	return params
}

func WriteErrorResponse(w http.ResponseWriter, err error, code int) {
	response := HttpResponse{
		Status:  "failed",
		Message: err.Error(),
	}

	w.Header().Set(ContentType, JsonHeader)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}
