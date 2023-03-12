package testgroup

import (
	"context"
	"github.com/codymj/go-service/foundation/web"
	"github.com/rs/zerolog"
	"net/http"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Logger *zerolog.Logger
}

// Test handler for development
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// ok status
	data := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	statusCode := http.StatusOK
	h.Logger.Info().Timestamp().
		Int("statusCode", statusCode).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remoteAddr", r.RemoteAddr).
		Msg("test")

	return web.Respond(ctx, w, data, statusCode)
}
