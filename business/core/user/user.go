package user

import (
	"context"
	"fmt"
	"github.com/codymj/go-service/business/data/store/user"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"time"
)

// Core manages the set of APIs for user access.
type Core struct {
	logger *zerolog.Logger
	user   user.Store
}

// NewCore constructs a core for user API access.
func NewCore(logger *zerolog.Logger, db *sqlx.DB) Core {
	return Core{
		logger: logger,
		user:   user.NewStore(logger, db),
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time) (user.User, error) {
	/*
	 * perform any pre-business operations
	 */

	// create user
	u, err := c.user.Create(ctx, nu, now)
	if err != nil {
		return user.User{}, fmt.Errorf("core.user.Create(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return u, nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, claims auth.Claims, userId string, uu user.UpdateUser, now time.Time) error {
	/*
	 * perform any pre-business operations
	 */

	// update user
	err := c.user.Update(ctx, claims, userId, uu, now)
	if err != nil {
		return fmt.Errorf("core.user.Update(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return nil
}

// Delete removes a user from the database.
func (c Core) Delete(ctx context.Context, claims auth.Claims, userId string) error {
	/*
	 * perform any pre-business operations
	 */

	// delete user
	if err := c.user.Delete(ctx, claims, userId); err != nil {
		return fmt.Errorf("core.user.Delete(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return nil
}

// Query retrives a list of existing users from the database.
func (c Core) Query(ctx context.Context, pageNum, rowsPerPg int) ([]user.User, error) {
	/*
	 * perform any pre-business operations
	 */

	// get user
	u, err := c.user.Query(ctx, pageNum, rowsPerPg)
	if err != nil {
		return []user.User{}, fmt.Errorf("core.user.Query(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return u, nil
}

// QueryById gets the specified user by ID from the database.
func (c Core) QueryById(ctx context.Context, claims auth.Claims, userId string) (user.User, error) {
	/*
	 * perform any pre-business operations
	 */

	// get user
	u, err := c.user.QueryById(ctx, claims, userId)
	if err != nil {
		return user.User{}, fmt.Errorf("core.user.QueryById(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return u, nil
}

// QueryByEmail gets the specified user by email from the database.
func (c Core) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (user.User, error) {
	/*
	 * perform any pre-business operations
	 */

	// get user
	u, err := c.user.QueryByEmail(ctx, claims, email)
	if err != nil {
		return user.User{}, fmt.Errorf("core.user.QueryByEmail(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return u, nil
}

// Authenticate finds a user by email and verifies their password.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	/*
	 * perform any pre-business operations
	 */

	// authenticate user
	claims, err := c.user.Authenticate(ctx, now, email, password)
	if err != nil {
		return auth.Claims{}, fmt.Errorf("core.user.Authenticate(): %w", err)
	}

	/*
	 * perform any post-business operations
	 */

	return claims, nil
}
