package user

import (
	"context"
	"go-service.codymj.io/internal/user/dao"
)

// service dependencies.
type service struct {
	userdao dao.Repository
}

// Service interface.
type Service interface {
	GetById(ctx context.Context, id int64) (dao.User, error)
	GetByParams(ctx context.Context, params map[string]string) ([]dao.User, error)
}

// New returns an initialized instance.
func New(dao dao.Repository) Service {
	return &service{
		dao: dao,
	}
}
