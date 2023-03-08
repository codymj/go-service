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
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	Registry *viper.Viper
)

type Web struct {
	APIHost         string
	DebugHost       string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// main service function =======================================================
func main() {
	// init logger
	logger := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	// run
	if err := run(&logger); err != nil {
		logger.Error().
			Err(err).
			Msg("error running api service")
		os.Exit(1)
	}
}

// run service =================================================================
func run(logger *zerolog.Logger) error {
	// set maxprocs
	opt := maxprocs.Logger(logger.Printf)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("%w", err)
	}
	logger.Info().
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Msg("startup")

	// set config parameters
	viper.AutomaticEnv()
	Registry = viper.GetViper()
	Registry.AddConfigPath(".")
	Registry.AddConfigPath("../../../..")
	Registry.SetConfigFile("conf.yml")
	err := Registry.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w \n", err)
	}

	cfg := struct {
		Web Web
	}{
		Web: Web{
			Registry.GetString("API_HOST"),
			Registry.GetString("DEBUG_HOST"),
			Registry.GetDuration("READ_TIMEOUT"),
			Registry.GetDuration("WRITE_TIMEOUT"),
			Registry.GetDuration("IDLE_TIMEOUT"),
			Registry.GetDuration("SHUTDOWN_TIMEOUT"),
		},
	}

	// start
	logger.Info().
		Msg("starting service")
	defer logger.Info().
		Msg("shutdown successful")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      nil,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	// shutdown
	select {
	case err = <-serverErrors:
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

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err = api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
