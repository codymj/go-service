package seed

import (
	"context"
	"database/sql"
	"sort"

	"github.com/ankorstore/yokai/config"
)

type UsersSeed struct {
	config *config.Config
}

func NewUsersSeed(cfg *config.Config) *UsersSeed {
	return &UsersSeed{
		config: cfg,
	}
}

func (s *UsersSeed) Name() string {
	return "users"
}

func (s *UsersSeed) Run(ctx context.Context, db *sql.DB) error {
	var txErr error

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	seedData := s.config.GetStringMap("config.seed.users")

	names := make([]string, 0, len(seedData))
	for name := range seedData {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		q := `
		INSERT INTO users(
			username,
			email,
			password,
			location,
			is_validated
		) VALUES(
			?, ?, ?, ?, ?
		)
		`

		username := seedData[name].(map[string]any)["username"].(string)
		email := seedData[name].(map[string]any)["email"].(string)
		password := seedData[name].(map[string]any)["password"].(string)
		var location sql.NullString
		if loc, found := seedData[name].(map[string]any)["location"].(string); found {
			location.Valid = true
			location.String = loc
		}
		is_validated := seedData[name].(map[string]any)["username"].(bool)

		_, txErr = tx.ExecContext(ctx, q, username, email, password, location, is_validated)
	}

	if txErr != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
