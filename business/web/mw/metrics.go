package mw

import (
	"context"
	"github.com/codymj/go-service/business/sys/metrics"
	"github.com/codymj/go-service/foundation/web"
	"net/http"
)

// Metrics updates program counters.
func Metrics() web.Middleware {
	// function to be executed
	m := func(handler web.Handler) web.Handler {
		// create the handler that will be connected to mw chain
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// add the metrics into the context for metrics gathering
			ctx = metrics.Set(ctx)

			// call next handler
			err := handler(ctx, w, r)

			// increment the request and goroutines counter
			metrics.AddRequest(ctx)
			metrics.AddGoroutine(ctx)

			// increment if there is an error
			if err != nil {
				metrics.AddError(ctx)
			}

			return err
		}

		return h
	}

	return m
}
