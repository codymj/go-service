package dao

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"strings"
	"time"
)

var (
	// A valid list of query parameters.
	validUserParams = []string{
		"userId",
		"firstName",
		"lastName",
		"email",
		"createdAt",
		"updatedAt",
	}

	// A map of query parameter names to database column names.
	paramToColumn = map[string]string{
		"userId":    "user_id",
		"firstName": "first_name",
		"lastName":  "last_name",
		"email":     "email",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
)

func getByEmailPasswordQuery() string {
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
	    1=1
	`
}

func (r *repository) GetByParams(ctx context.Context, params map[string]string) ([]User, error) {
	// Build "where" clause for query.
	query := getByEmailPasswordQuery()
	whereClause, vals := buildWhereClause(params)
	if !strings.EqualFold("", whereClause) {
		query = strings.Replace(query, "1=1", whereClause, 1)
	}

	// Execute query.
	rows, err := r.db.DB.QueryContext(ctx, query, vals...)
	if err != nil {
		return nil, errors.Wrap(err, ErrQueryingDatabase.Error())
	}
	defer Close(&err, io.Closer(rows))

	// Parse result.
	users := make([]User, 0)
	for rows.Next() {
		var id int64
		var firstName string
		var lastName string
		var email string
		var createdAt time.Time
		var updatedAt time.Time

		err = rows.Scan(
			&id, &firstName, &lastName, &email, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, ErrParsingRow.Error())
		}

		user := User{
			UserId:    id,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		users = append(users, user)
	}

	return users, nil
}
