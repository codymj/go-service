package commands

import (
	"context"
	"github.com/codymj/go-service/business/data/schema"
	"github.com/codymj/go-service/business/sys/database"
	"time"
)

// Seed seeds the database with test data
func Seed(config database.Config) error {
	// open db connection
	db, err := database.Open(config)
	if err != nil {
		return err
	}
	defer db.Close()

	// set a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// migrate
	if err = schema.Seed(ctx, db); err != nil {
		return err
	}

	return nil
}
