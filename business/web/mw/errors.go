package mw

import (
	"context"
	"github.com/codymj/go-service/business/sys/validate"
	"github.com/codymj/go-service/foundation/web"
	"github.com/rs/zerolog"
	"net/http"
)

// Errors handles errors coming out of the call chain. It detects normal app
// errors which are used to respond to the client in a uniform way. Unexpected
// errors (5xx) are logged.
func Errors(logger *zerolog.Logger) web.Middleware {
	// function to be executed
	m := func(handler web.Handler) web.Handler {
		// create the handler that will be connected to mw chain
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// if the context is missing this value, shutdown service
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			// run the next handler and catch any propagated error
			if err = handler(ctx, w, r); err != nil {
				// log error
				logger.Error().Timestamp().
					Str("traceId", v.TraceId).
					Err(err).
					Msg("ERROR")

				// build error response
				var e validate.ErrorResponse
				var status int
				switch act := validate.Cause(err).(type) {
				case validate.FieldErrors:
					e = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: act.Error(),
					}
					status = http.StatusBadRequest
				case *validate.RequestError:
					e = validate.ErrorResponse{
						Error: act.Error(),
					}
					status = act.Status
				default:
					e = validate.ErrorResponse{
						Error: "internal error",
					}
					status = http.StatusInternalServerError
				}

				// respond with the error back to client
				if errr := web.Respond(ctx, w, e, status); errr != nil {
					return errr
				}

				// if we receive the shutdown err we need to return it back to
				// the base handler to shut down the service
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// error has been handled, no need to propagate it
			return nil
		}

		return h
	}

	return m
}
