package handlers

import (
	"expvar"
	"github.com/codymj/go-service/app/services/api/handlers/debug/checkgroup"
	"github.com/codymj/go-service/app/services/api/handlers/v1/testgroup"
	"github.com/codymj/go-service/business/web/mw"
	"github.com/codymj/go-service/foundation/web"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/pprof"
	"os"
)

// DebugStdLibMux registers all the debug routes from the standard library into
// a new mux bypassing the use of the DefaultServerMux
func DebugStdLibMux() *http.ServeMux {
	// build mux
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service
func DebugMux(build string, logger *zerolog.Logger) http.Handler {
	// register debug check endpoints
	cgh := checkgroup.Handlers{
		Build:  build,
		Logger: logger,
	}

	// build mux
	mux := DebugStdLibMux()
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// ApiMuxConfig contains all the mandatory systems required by handlers.
type ApiMuxConfig struct {
	Shutdown chan os.Signal
	Logger   *zerolog.Logger
}

// ApiMux constructs an http.Handler with all application routes defined.
func ApiMux(cfg ApiMuxConfig) *web.App {
	// create new web application
	app := web.NewApp(
		cfg.Shutdown,
		mw.Logger(cfg.Logger),
	)

	// bind v1 routes
	v1(app, cfg)

	return app
}

// v1 binds all the v1 routes
func v1(app *web.App, cfg ApiMuxConfig) {
	const version = "v1"

	// create v1 test handler
	tgh := testgroup.Handlers{
		Logger: cfg.Logger,
	}

	// register test routes
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
}
