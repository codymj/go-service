package metrics

import (
	"context"
	"expvar"
)

// This holds the single instance of the metrics value needed for collecting
// metrics. The expvar package is already based on a singleton for the different
// metrics that are registered with the package.
var m *metrics

// metrics represents the set of metrics we gather. These fields are safe to be
// accessed concurrently thanks to expvar.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constructs the metrics value that will be used to capture metrics.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// ctxKeyMetric represents the type of value for the context key.
type ctxKey int

// key is ctxKey
const key ctxKey = 1

// Set sets the metrics data into the context
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutine increments the goroutines metric by 1.
func AddGoroutine(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Add(1)
		}
	}
}

// AddRequest increments the requests metric by 1.
func AddRequest(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddError increments the errors metric by 1.
func AddError(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanic increments the panics metric by 1.
func AddPanic(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}
