package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/sql/healthcheck"
	"go-service.codymj.io/internal/repository"
	"go-service.codymj.io/internal/service"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// Services:
		fx.Provide(repository.NewUserRepository, service.NewUserService),

		// Metrics:
		fxmetrics.AsMetricsCollector(service.UserServiceCounter),

		// Probes:
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
	)
}
