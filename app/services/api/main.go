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
	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"
)

var build = "dev"

// main service function
func main() {
	// init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// run
	if err := run(); err != nil {
		log.Error().
			Err(err).
			Msg("error running api service")
		os.Exit(1)
	}
}

// run service
func run() error {
	// maxproc
	opt := maxprocs.Logger(log.Printf)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Info().
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Msg("startup")

	// start
	log.Info().
		Str("version", build).
		Msg("starting service")
	defer log.Info().
		Msg("shutdown successful")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         "127.0.0.1:3333",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info().
			Str("status", "starting").
			Msg("starting service")
		serverErrors <- api.ListenAndServe()
	}()

	// shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info().
			Str("status", "shutdown").
			Any("signal", sig).
			Msg("shutting down")
		defer log.Info().
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
