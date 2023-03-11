package testgroup

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Logger *zerolog.Logger
}

// Test handler for development
func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	// ok status
	data := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	_ = json.NewEncoder(w).Encode(data)

	statusCode := http.StatusOK
	h.Logger.Info().Timestamp().
		Int("statusCode", statusCode).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remoteAddr", r.RemoteAddr).
		Msg("test")
}
