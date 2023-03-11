package checkgroup

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Build  string
	Logger *zerolog.Logger
}

// Readiness checks if the database is ready and if not, returns 500
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	// ok status
	data := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	// build response
	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		// log error
		h.Logger.Error().Timestamp().
			Err(err).
			Msg("readiness failed")
	}

	// log info
	h.Logger.Info().Timestamp().
		Int("statusCode", statusCode).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remoteAddr", r.RemoteAddr).
		Msg("readiness")
}

// Liveness returns status info pertaining to the service
func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	// host information
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "down"
	}

	// ok status
	data := struct {
		Status   string `json:"status"`
		Build    string `json:"build"`
		Hostname string `json:"hostname"`
		Host     string `json:"host"`
	}{
		Status:   "ok",
		Build:    h.Build,
		Hostname: hostname,
		Host:     r.RemoteAddr,
	}

	// build response
	statusCode := http.StatusOK
	if err = response(w, statusCode, data); err != nil {
		// log error
		h.Logger.Error().Timestamp().
			Err(err).
			Msg("liveness failed")
	}

	// log info
	h.Logger.Info().Timestamp().
		Int("statusCode", statusCode).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remoteAddr", r.RemoteAddr).
		Msg("liveness")
}

// util functions ==============================================================

func response(w http.ResponseWriter, statusCode int, data any) error {
	// convert response to json
	jsn, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// set content type and headers
	w.Header().Set("Content-Type", "application/json")

	// write status code
	w.WriteHeader(statusCode)

	// send
	if _, err = w.Write(jsn); err != nil {
		return err
	}

	return nil
}
