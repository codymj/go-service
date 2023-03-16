package mw

import (
	"context"
	"fmt"
	"github.com/codymj/go-service/business/sys/metrics"
	"github.com/codymj/go-service/foundation/web"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panics and converts the panic to an error, so it is
// reported in Metrics and handled in Errors.
func Panics() web.Middleware {
	// function to be executed
	m := func(handler web.Handler) web.Handler {
		// create the handler that will be connected to mw chain
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// defer function to recover from panic, setting error
			defer func() {
				if rec := recover(); rec != nil {
					// stack trace
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE %+v", rec, string(trace))

					// update metrics
					metrics.AddPanic(ctx)
				}
			}()

			// call the next handler and set return value in the err variable
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
