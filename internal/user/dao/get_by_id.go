package dao

import (
	"context"
	"strconv"
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
}
