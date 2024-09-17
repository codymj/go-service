package dao

import (
	"context"
	"go-service.codymj.io/internal/database"
)

type repository struct {
	db *database.Connection
	ps password.Service
}

type Repository interface {
	GetById(ctx context.Context, id int64) (User, error)
	GetByParams(ctx context.Context, params map[string]string) ([]User, error)
}

func New(db *database.Connection, ps password.Service) Repository {
	return &repository{
		db: db,
		ps: ps,
	}
}
