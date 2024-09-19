package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Version for this application.
const version = "0.0.1"

// Application configuration properties.
type config struct {
	port int
	env  string
}

// Application dependencies.
type app struct {
	config config
	logger *slog.Logger
}

// Application's main function.
func main() {
	var cfg config

	// Parse commandline flags.
	flag.IntVar(&cfg.port, "port", 8080, "Application listening port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stg|prd)")
	flag.Parse()

	// Initialize logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize application.
	app := &app{
		config: cfg,
		logger: logger,
	}

	// Initialize web server.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start web server.
	logger.Info("Starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
