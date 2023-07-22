package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/codymj/go-service/business/data/schema"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/codymj/go-service/foundation/docker"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"io"
	"os"
	"testing"
	"time"
)

// success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// DbContainer provides configuration for a container to run.
type DbContainer struct {
	Image string
	Port  string
	Args  []string
}

// NewUnit creates a seeded test database in a Docker container.
func NewUnit(t *testing.T, dbc DbContainer) (*zerolog.Logger, *sqlx.DB, func()) {
	// redirect stdout to a pipe for logging
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	// start container
	c := docker.StartContainer(t, dbc.Image, dbc.Port, dbc.Args...)

	// connect to database
	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTls: true,
	})
	if err != nil {
		t.Fatalf("database.Open(): %v", err)
	}
	t.Log("waiting for database to be ready...")

	// run migration and seed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = schema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.Id)
		docker.StopContainer(t, c.Id)
		t.Fatalf("schema.Migrate(): %v", err)
	}
	if err = schema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.Id)
		docker.StopContainer(t, c.Id)
		t.Fatalf("schema.Seed(): %v", err)
	}

	// log
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	// teardown
	teardown := func() {
		t.Helper()
		db.Close()
		docker.StopContainer(t, c.Id)

		_ = os.Stdout.Sync()

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("************************* LOGS *************************")
		fmt.Print(buf.String())
		fmt.Println("************************* LOGS *************************")
	}

	return &logger, db, teardown
}

// StrPtr is a helper to return *string for a string.
func StrPtr(s string) *string {
	return &s
}

// IntPtr is a helper to return *int for an int.
func IntPtr(i int) *int {
	return &i
}
