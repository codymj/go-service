package handlers

import (
	"expvar"
	"github.com/codymj/go-service/app/services/api/handlers/debug/checkgroup"
	"github.com/codymj/go-service/app/services/api/handlers/v1/testgroup"
	"github.com/codymj/go-service/app/services/api/handlers/v1/usergroup"
	usercore "github.com/codymj/go-service/business/core/user"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/codymj/go-service/business/web/mw"
	"github.com/codymj/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
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
func DebugMux(build string, logger *zerolog.Logger, db *sqlx.DB) http.Handler {
	// register debug check endpoints
	cgh := checkgroup.Handlers{
		Build:  build,
		Logger: logger,
		DB:     db,
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
	Auth     *auth.Auth
	DB       *sqlx.DB
}

// ApiMux constructs an http.Handler with all application routes defined.
func ApiMux(cfg ApiMuxConfig) *web.App {
	// create new web application
	app := web.NewApp(
		cfg.Shutdown,
		mw.Logger(cfg.Logger),
		mw.Errors(cfg.Logger),
		mw.Metrics(),
		mw.Panics(),
	)

	// bind v1 routes
	v1(app, cfg)

	return app
}

// v1 binds all the v1 routes.
func v1(app *web.App, cfg ApiMuxConfig) {
	// set API version
	const version = "v1"

	// test handler and routes
	tgh := testgroup.Handlers{
		Logger: cfg.Logger,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, mw.Authenticate(cfg.Auth), mw.Authorize("admin"))

	// user handler and routes
	ugh := usergroup.Handlers{
		Auth: cfg.Auth,
		User: usercore.NewCore(cfg.Logger, cfg.DB),
	}
	app.Handle(
		http.MethodGet,
		version,
		"/users/token",
		ugh.Token)
	app.Handle(
		http.MethodGet,
		version,
		"/users/:page/:rows",
		ugh.Query,
		mw.Authenticate(cfg.Auth),
		mw.Authorize(auth.RoleAdmin))
	app.Handle(
		http.MethodGet,
		version,
		"/users/:id",
		ugh.QueryById,
		mw.Authenticate(cfg.Auth))
	app.Handle(
		http.MethodPost,
		version,
		"/users",
		ugh.Create,
		mw.Authenticate(cfg.Auth),
		mw.Authorize(auth.RoleAdmin))
	app.Handle(
		http.MethodPut,
		version,
		"/users/:id",
		ugh.Update,
		mw.Authenticate(cfg.Auth),
		mw.Authorize(auth.RoleAdmin))
	app.Handle(
		http.MethodDelete,
		version,
		"/users/:id",
		ugh.Delete,
		mw.Authenticate(cfg.Auth),
		mw.Authorize(auth.RoleAdmin))
}
