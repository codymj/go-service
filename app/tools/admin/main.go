package main

import (
	"flag"
	"fmt"
	"github.com/codymj/go-service/app/tools/admin/commands"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/codymj/go-service/foundation/keystore"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
)

var (
	Registry *viper.Viper
)

type config struct {
	AuthCfg keystore.AuthCfg
	DbCfg   database.Config
}

func main() {
	// init logger
	logger := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel)

	// run
	err := run(&logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// run service =================================================================
func run(logger *zerolog.Logger) error {
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
		AuthCfg keystore.AuthCfg
		DbCfg   database.Config
	}{
		AuthCfg: keystore.AuthCfg{
			KeysFolder: "zarf/keys/",
			ActiveKid:  "1b24502a-4781-47cb-99c2-3403c23bedac",
		},
		DbCfg: database.Config{
			User:         Registry.GetString("DB_USER"),
			Password:     Registry.GetString("DB_PASSWORD"),
			Host:         Registry.GetString("DB_HOST"),
			Name:         Registry.GetString("DB_NAME"),
			MaxIdleConns: Registry.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns: Registry.GetInt("DB_MAX_OPEN_CONNS"),
			DisableTls:   Registry.GetBool("DB_TLS_DISABLED"),
		},
	}

	// flags
	flag.Parse()

	return processCommands(logger, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(logger *zerolog.Logger, cfg config) error {
	switch os.Args[1] {
	case "genkeys":
		if err := commands.GenKeys(); err != nil {
			return fmt.Errorf("commands.GenKeys(): %w", err)
		}
	case "gentoken":
		if err := commands.GenToken(logger, cfg.DbCfg, cfg.AuthCfg); err != nil {
			return fmt.Errorf("commands.GenToken(): %w", err)
		}
	case "migrate":
		if err := commands.Migrate(cfg.DbCfg); err != nil {
			return fmt.Errorf("commands.Migrate(): %w", err)
		}
	case "seed":
		if err := commands.Seed(cfg.DbCfg); err != nil {
			return fmt.Errorf("commands.Seed(): %w", err)
		}
	}

	return nil
}
