package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"os"
	"strings"
	"syscall"
)

// Handler is a type that handles an http request within our own framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and configures the context object
// for each of the http handlers.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given http method and path.
func (a *App) Handle(method, group, path string, handler Handler, mw ...Middleware) {
	// wrap handler-specific middleware
	handler = wrapMiddleware(mw, handler)

	// wrap the application-specific middleware
	handler = wrapMiddleware(a.mw, handler)

	// execute request
	h := func(w http.ResponseWriter, r *http.Request) {
		// pre-handler processing

		// call wrapped handler functions
		if err := handler(r.Context(), w, r); err != nil {
			// handle error
			return
		}

		// post-handler processing
	}

	// construct final path
	finalPath := path
	if !strings.EqualFold("", group) {
		finalPath = "/" + group + path
	}

	a.ContextMux.Handle(method, finalPath, h)
}
