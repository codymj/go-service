package schema

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/ardanlabs/darwin"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/delete.sql
	deleteDoc string
)

// Migrate attempts to bring the schema for database up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	// check db status
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("database.StatusCheck(): %w", err)
	}

	// construct darwin driver
	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return fmt.Errorf("darwin.NewGenericDriver(): %w", err)
	}

	// build darwin struct
	drwn := darwin.New(driver, darwin.ParseMigrations(schemaDoc))

	// migrate
	err = drwn.Migrate()
	if err != nil {
		return fmt.Errorf("drwn.Migrate(): %w", err)
	}

	return nil
}

// DeleteAll runs the set of DELETE queries against database. The queries are ran
// in a transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	// build transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin(): %w", err)
	}

	// execute transaction
	if _, err0 := tx.Exec(deleteDoc); err != nil {
		if err1 := tx.Rollback(); err != nil {
			return fmt.Errorf("tx.Rollback(): %w", err1)
		}
		return fmt.Errorf("tx.Exec() delete: %w", err0)
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}
