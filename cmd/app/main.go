package main

import (
	"fmt"
	"go-service.codymj.io/cmd/app/config"
	"go-service.codymj.io/cmd/app/util"
	"go-service.codymj.io/internal/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	err := start()
	if err != nil {
		panic(fmt.Errorf("Fatal error starting service: %w\n", err))
	}
}

func start() error {
	// Initialize application.
	config.Set()
	port := config.Registry.GetInt("SERVER_PORT")
	readTimeoutStr := config.Registry.GetString("SERVER.READ_TIMEOUT")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		log.Fatal().Msg("Set a valid server read timeout duration.")
		return err
	}
	writeTimeoutStr := config.Registry.GetString("SERVER.WRITE_TIMEOUT")
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		log.Fatal().Msg("Set a valid server write timeout duration.")
		return err
	}
	idleTimeoutStr := config.Registry.GetString("SERVER.IDLE_TIMEOUT")
	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		log.Fatal().Msg("Set a valid server idle timeout duration.")
		return err
	}

	// Initialize database.
	db, err := database.Configure(config.GetDatabaseConfig())
	if err != nil {
		log.Fatal().Msg("Error initializing database connection.")
		return err
	}

	// Initialize utility services.
	validateSvc := config.NewValidateService()
	passwordSvc := config.NewPasswordService(config.GetPasswordConfig())

	// Initialize repositories.
	userRepo := config.NewUserRepository(db, passwordSvc)

	// Initialize business services.
	userSvc := config.NewUserService(userRepo)

	// Initialize routes.
	routeServices := util.Services{
		ValidatorService: validateSvc,
		UserService:      userSvc,
	}
	router := route.NewRouter()
	err = router.Setup(routeServices)
	if err != nil {
		log.Fatal().Msg("Error initializing HTTP routes.")
		return err
	}

	// Start application.
	log.Info().Msg(fmt.Sprintf("Service running on port %d", port))
	srv := &http.Server{
		Handler:      router.Router,
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal().Msg("Error starting application")
		return err
	}

	return nil
}
