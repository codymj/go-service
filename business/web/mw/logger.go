package mw

import (
	"context"
	"github.com/codymj/go-service/foundation/web"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

// Logger middleware
func Logger(logger *zerolog.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// if the context is missing this value, shutdown service
			v, err := web.GetValues(ctx)
			if err != nil {
				return err // todo: handle shutdown
			}

			logger.Info().Timestamp().
				Str("traceId", v.TraceId).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remoteAddr", r.RemoteAddr).
				Msg("request started")

			err = handler(ctx, w, r)

			logger.Info().Timestamp().
				Str("traceId", v.TraceId).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remoteAddr", r.RemoteAddr).
				Int("statusCode", v.StatusCode).
				Dur("since", time.Since(v.Now)).
				Msg("request completed")

			return err
		}

		return h
	}
}
