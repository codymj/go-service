package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go-service.codymj.io/cmd/app/router/user"
	"go-service.codymj.io/internal/database"
	"go-service.codymj.io/internal/password"
	userdao "go-service.codymj.io/internal/user/dao"
	"time"
)

var (
	Registry *viper.Viper
)

func Set() {
	viper.AutomaticEnv()

	Registry = viper.GetViper()
	Registry.AddConfigPath(".")
	Registry.AddConfigPath("../..")
	Registry.SetConfigFile("configuration.yml")
	err := Registry.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error with file: %w\n", err))
	}

	SetLoggerParams()
}

func SetLoggerParams() {
	zerolog.TimeFieldFormat = time.RFC3339

	loggerLevel := Registry.GetString("LOGGER.LEVEL")
	switch loggerLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		log.Warn().Msg("No log level set, defaulting to info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func GetDatabaseConfig() *database.Config {
	connMaxLifetimeStr := Registry.GetString("POSTGRES.CONN_MAX_LIFETIME")
	connMaxLifetime, err := time.ParseDuration(connMaxLifetimeStr)
	if err != nil {
		log.Warn().Msg("Invalid POSTGRES.CONN_MAX_LIFETIME, setting to 3m")
		connMaxLifetime, _ = time.ParseDuration("3m")
	}
	connMaxIdleTimeStr := Registry.GetString("POSTGRES.CONN_MAX_IDLE_TIME")
	connMaxIdleTime, err := time.ParseDuration(connMaxIdleTimeStr)
	if err != nil {
		log.Warn().Msg("Invalid POSTGRES.CONN_MAX_IDLE_TIME, setting to 15m")
		connMaxLifetime, _ = time.ParseDuration("15m")
	}

	return &database.Config{
		User:            Registry.GetString("POSTGRES.USER"),
		Password:        Registry.GetString("POSTGRES.PASSWORD"),
		Name:            Registry.GetString("POSTGRES.NAME"),
		Host:            Registry.GetString("POSTGRES.HOST"),
		Port:            Registry.GetInt("POSTGRES.PORT"),
		ConnMaxLifetime: connMaxLifetime,
		MaxOpenConns:    Registry.GetInt("POSTGRES.MAX_OPEN_CONNS"),
		MaxIdleConns:    Registry.GetInt("POSTGRES.MAX_IDLE_CONNS"),
		ConnMaxIdleTime: connMaxIdleTime,
	}
}

func GetPasswordConfig() *password.Config {
	return &password.Config{
		Time:      1,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}
}

func NewUserRepository(db *database.Connection, ps password.Service) userdao.Repository {
	return userdao.New(db, ps)
}

func NewValidateService() validate.Service {
	return validate.New()
}

func NewPasswordService(cfg *password.Config) password.Service {
	return password.New(cfg)
}

func NewUserService(ur userdao.Repository) user.Service {
	return user.New(ur)
}
