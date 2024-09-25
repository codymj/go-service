package repository

import (
	"context"
	"database/sql"
	"sync"

	"github.com/Masterminds/squirrel"
	"go-service.codymj.io/internal/model"
)

// UserRepository is the repository to handle the model.User model database interractions.
type UserRepository struct {
	mutex sync.Mutex
	db    *sql.DB
}

// NewUserRepository returns a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// UserRepositoryFindAllParams contains parameters for FindAll.
type UserRepositoryFindAllParams struct {
	Username sql.NullString
	Email    sql.NullString
}

// FindAll finds all users and filters by parameters.
func (r *UserRepository) FindAll(
	ctx context.Context,
	params UserRepositoryFindAllParams,
) ([]model.User, error) {
	// Build query.
	qb := squirrel.
		Select(
			"id",
			"username",
			"email",
			"location",
			"is_validated",
			"created_at",
			"modified_at",
		).
		From(
			"users",
		)

	if params.Username.Valid {
		qb = qb.Where(squirrel.Eq{"username": params.Username.String})
	}
	if params.Email.Valid {
		qb = qb.Where(squirrel.Eq{"email": params.Email.String})
	}
	qb = qb.OrderBy("id")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	// Query database.
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// Build response.
	var users []model.User
	for rows.Next() {
		var user model.User
		if err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Location,
			&user.CreatedAt,
			&user.ModifiedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
