package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/automaxprocs/maxprocs"
)

var build = "dev"

// main service function
func main() {
	// init logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()

	// run
	if err := run(&logger); err != nil {
		logger.Error().
			Err(err).
			Msg("error running api service")
		os.Exit(1)
	}
}

// run service
func run(logger *zerolog.Logger) error {
	// maxproc
	opt := maxprocs.Logger(logger.Printf)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	logger.Info().
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Msg("startup")

	// start
	logger.Info().
		Str("version", build).
		Msg("starting service")
	defer logger.Info().
		Msg("shutdown successful")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         "127.0.0.1:3333",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info().
			Str("status", "starting").
			Msg("starting service")
		serverErrors <- api.ListenAndServe()
	}()

	// shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info().
			Str("status", "shutdown").
			Any("signal", sig).
			Msg("shutting down")
		defer logger.Info().
			Str("status", "stopped").
			Any("signal", sig).
			Msg("shutdown complete")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
