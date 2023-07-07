package checkgroup

import (
	"context"
	"encoding/json"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Build  string
	Logger *zerolog.Logger
	DB     *sqlx.DB
}

// Readiness checks if the database is ready and if not, returns 500
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	// give one second for readiness to pass all tests
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	// check db
	status := "ok"
	msg := "readiness"
	statusCode := http.StatusOK
	if err := database.StatusCheck(ctx, h.DB); err != nil {
		status = "error"
		msg = err.Error()
		statusCode = http.StatusInternalServerError
	}

	// build response
	data := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  status,
		Message: msg,
	}
	if err := response(w, statusCode, data); err != nil {
		// log error
		h.Logger.Error().Timestamp().
			Err(err).
			Msg("readiness failed")
	}

	// log info
	h.Logger.Info().Timestamp().
		Str("status", status).
		Int("statusCode", statusCode).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remoteAddr", r.RemoteAddr).
		Msg(msg)
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
