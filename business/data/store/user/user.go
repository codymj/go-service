package user

import (
	"context"
	"fmt"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/codymj/go-service/business/sys/validate"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Store manages the set of APIs for user access.
type Store struct {
	logger *zerolog.Logger
	db     *sqlx.DB
}

// NewStore constructs a user store for API access.
func NewStore(logger *zerolog.Logger, db *sqlx.DB) Store {
	return Store{
		logger: logger,
		db:     db,
	}
}

// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	// validate new user request
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validate.Check(): %w", err)
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("bcrypt.GenerateFromPassword(): %w", err)
	}

	// create new user
	u := User{
		Id:           validate.GenerateId(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		Created:      now,
		Modified:     now,
	}

	// query to insert new user
	const q = `
	INSERT INTO users (
		user_id, name, email, password_hash, roles, created, modified)
	VALUES (
		:user_id, :name, :email, :password_hash, :roles, :created, :modified)
	`

	// insert new user into database
	if err = database.NamedExecContext(ctx, s.logger, s.db, q, u); err != nil {
		return User{}, fmt.Errorf("database.NamedExecContext(): %w", err)
	}

	return u, nil
}

// Update replaces a user in the database.
func (s Store) Update(ctx context.Context, claims auth.Claims, userId string, uu UpdateUser, now time.Time) error {
	// validate update user request
	if err := validate.CheckId(userId); err != nil {
		return database.ErrInvalidId
	}
	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validate.Check(): %w", err)
	}

	// get user from database and update
	u, err := s.QueryById(ctx, claims, userId)
	if err != nil {
		return fmt.Errorf("s.QueryById(): %w", err)
	}
	if uu.Name != nil {
		u.Name = *uu.Name
	}
	if uu.Email != nil {
		u.Email = *uu.Email
	}
	if uu.Roles != nil {
		u.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("bcrypt.GeneratePassword(): %w", err)
		}
		u.PasswordHash = pw
	}
	u.Modified = now

	// query to update user
	const q = `
	UPDATE
		users
	SET
	    "name" = :name,
	    "email" = :email,
	    "roles" = :roles,
	    "password_hash" = :password_hash,
	    "modified" = :modified
	WHERE
	    user_id = :user_id
	`

	// update user
	if err = database.NamedExecContext(ctx, s.logger, s.db, q, u); err != nil {
		return fmt.Errorf("database.NamedExecContext(): %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (s Store) Delete(ctx context.Context, claims auth.Claims, userId string) error {
	// validate
	if err := validate.CheckId(userId); err != nil {
		return database.ErrInvalidId
	}

	// check authorization
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userId {
		return database.ErrForbidden
	}

	// data to delete user
	data := struct {
		UserId string `db:"user_id"`
	}{
		UserId: userId,
	}

	// query to delete user
	const q = `
	DELETE FROM
		users
	WHERE
	    user_id = :user_id
	`

	// delete user
	if err := database.NamedExecContext(ctx, s.logger, s.db, q, data); err != nil {
		return fmt.Errorf("database.NamedExecContext(): %w", err)
	}

	return nil
}

// Query retrives a list of existing users from the database.
func (s Store) Query(ctx context.Context, pageNum, rowsPerPg int) ([]User, error) {
	// data for query
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNum - 1) * rowsPerPg,
		RowsPerPage: rowsPerPg,
	}

	// query to get users
	const q = `
	SELECT
		*
	FROM
	    users
	ORDER BY
	    user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY
	`

	// get users
	var users []User
	if err := database.NamedQuerySlice(ctx, s.logger, s.db, q, data, &users); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("database.NamedQuerySlice(): %w", err)
	}

	return users, nil
}

// QueryById gets the specified user by ID from the database.
func (s Store) QueryById(ctx context.Context, claims auth.Claims, userId string) (User, error) {
	// validate
	if err := validate.CheckId(userId); err != nil {
		return User{}, database.ErrInvalidId
	}

	// check authorization
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userId {
		return User{}, database.ErrForbidden
	}

	// data to delete user
	data := struct {
		UserId string `db:"user_id"`
	}{
		UserId: userId,
	}

	// query to get user
	const q = `
	SELECT
		*
	FROM
	    users
	WHERE
	    user_id = :user_id
	`

	// get user
	var u User
	if err := database.NamedQueryStruct(ctx, s.logger, s.db, q, data, &u); err != nil {
		if err == database.ErrNotFound {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("database.NamedQueryStruct(): %w", err)
	}

	return u, nil
}

// QueryByEmail gets the specified user by email from the database.
func (s Store) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (User, error) {
	// todo: validate email

	// data to delete user
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	// query to get user
	const q = `
	SELECT
		*
	FROM
	    users
	WHERE
	    email = :email
	`

	// get user
	var u User
	if err := database.NamedQueryStruct(ctx, s.logger, s.db, q, data, &u); err != nil {
		if err == database.ErrNotFound {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("database.NamedQueryStruct(): %w", err)
	}

	// check authorization
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != u.Id {
		return User{}, database.ErrForbidden
	}

	return u, nil
}

// Authenticate finds a user by email and verifies their password.
func (s Store) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	// todo: validate email

	// data to delete user
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	// query to get user
	const q = `
	SELECT
		*
	FROM
	    users
	WHERE
	    email = :email
	`

	// get user
	var u User
	if err := database.NamedQueryStruct(ctx, s.logger, s.db, q, data, &u); err != nil {
		if err == database.ErrNotFound {
			return auth.Claims{}, database.ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("database.NamedQueryStruct(): %w", err)
	}

	// compare provided password to saved hash
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, database.ErrAuthFailure
	}

	// request is valid
	return auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "go-service",
			Subject: u.Id,
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(time.Hour).UTC(),
			},
			IssuedAt: &jwt.NumericDate{
				Time: now.UTC(),
			},
		},
		Roles: u.Roles,
	}, nil
}
