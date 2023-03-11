package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/codymj/go-service/app/services/api/handlers"
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

type App struct {
	BuildVersion string
}

type Web struct {
	ApiHost         string
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
		Level(zerolog.InfoLevel)

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

	// set config parameters from conf.yml
	viper.AutomaticEnv()
	Registry = viper.GetViper()
	Registry.AddConfigPath(".")
	Registry.AddConfigPath("../../../..")
	Registry.SetConfigFile("conf.yml")
	err := Registry.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w \n", err)
	}

	// initialize config
	cfg := struct {
		App App
		Web Web
	}{
		App: App{
			Registry.GetString("BUILD_VERSION"),
		},
		Web: Web{
			Registry.GetString("API_HOST"),
			Registry.GetString("DEBUG_HOST"),
			Registry.GetDuration("READ_TIMEOUT"),
			Registry.GetDuration("WRITE_TIMEOUT"),
			Registry.GetDuration("IDLE_TIMEOUT"),
			Registry.GetDuration("SHUTDOWN_TIMEOUT"),
		},
	}
	expvar.NewString("build").Set(cfg.App.BuildVersion)

	// start debug service
	logger.Info().Timestamp().
		Str("status", "started").
		Str("host", cfg.Web.DebugHost).
		Msg("debug router started")

	// start debug service
	debugMux := handlers.DebugMux(cfg.App.BuildVersion, logger)
	go func() {
		err = http.ListenAndServe(cfg.Web.DebugHost, debugMux)
		if err != nil {
			logger.Error().Timestamp().
				Str("status", "shutdown").
				Str("host", cfg.Web.DebugHost).
				Err(err).
				Msg("debug service shutting down")
		}
	}()

	// start api service
	logger.Info().Timestamp().
		Str("status", "started").
		Str("build", cfg.App.BuildVersion).
		Str("host", cfg.Web.ApiHost).
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Msg("service started")

	// buffered channel to listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// construct api mux
	apiMux := handlers.ApiMux(handlers.ApiMuxConfig{
		Shutdown: shutdown,
		Logger:   logger,
	})

	// init http server
	api := http.Server{
		Addr:         cfg.Web.ApiHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	// buffered channel to listen for server errors
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	// block main and wait for shutdown or server error
	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Info().Timestamp().
			Str("status", "shutdown").
			Str("host", cfg.Web.ApiHost).
			Any("signal", sig.String()).
			Msg("service shutting down")

		// give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// ask listener to shut down and shed the load
		if err = api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	logger.Info().Timestamp().
		Str("status", "stopped").
		Str("host", cfg.Web.ApiHost).
		Msg("service stopped")

	return nil
}
