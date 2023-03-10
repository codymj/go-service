package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
)

// DebugStdLibMux registers all the debug routes from the standard library into
// a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a depencency could inject a
// handler into our server without permission.
func DebugStdLibMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
