package dao

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func getByIdQuery() string {
	return `
	select
		user_id,
		first_name,
		last_name,
		email,
		created_at,
		updated_at
	from
	    users
	where
	    user_id = $1
	`
}

func (r *repository) GetById(ctx context.Context, id int64) (User, error) {
	// Execute query.
	row := r.db.DB.QueryRowContext(ctx, getByIdQuery(), strconv.Itoa(int(id)))

	// Parse result.
	var firstName string
	var lastName string
	var email string
	var createdAt time.Time
	var updatedAt time.Time

	err := row.Scan(
		&id, &firstName, &lastName, &email, &createdAt, &updatedAt,
	)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return User{}, nil
	} else if err != nil {
		return User{}, errors.Wrap(err, ErrParsingRow.Error())
	}

	user := User{
		UserId:    id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}
