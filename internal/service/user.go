package service

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/prometheus/client_golang/prometheus"
	"go-service.codymj.io/internal/model"
	"go-service.codymj.io/internal/repository"
)

// UserServiceCounter is a counter for the operation on users.
var UserServiceCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "user_service_operations_total",
		Help: "Number of operations on the UserService",
	},
	[]string{
		"operation",
	},
)

// UserService is the service to manage users.
type UserService struct {
	config     *config.Config
	repository *repository.UserRepository
}

// NewUserService returns a new UserService.
func NewUserService(cfg *config.Config, repo *repository.UserRepository) *UserService {
	return &UserService{
		config:     cfg,
		repository: repo,
	}
}

// List returns a list of all users, filtered by optional query parameters.
func (s *UserService) List(
	ctx context.Context,
	params repository.UserRepositoryFindAllParams,
) ([]model.User, error) {
	UserServiceCounter.WithLabelValues("list").Inc()

	return s.repository.FindAll(ctx, params)
}
