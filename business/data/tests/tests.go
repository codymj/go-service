package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/codymj/go-service/business/data/schema"
	"github.com/codymj/go-service/business/data/store/user"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/codymj/go-service/foundation/docker"
	"github.com/codymj/go-service/foundation/keystore"
	"github.com/google/uuid"
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

// Test owns state for running and shutting down tests.
type Test struct {
	Db       *sqlx.DB
	Logger   *zerolog.Logger
	Auth     *auth.Auth
	Teardown func()
	t        *testing.T
}

// Token generates an authenticated token for a user.
func (t *Test) Token(email, password string) string {
	t.t.Log("generating token for test...")

	// init store
	store := user.NewStore(t.Logger, t.Db)
	claims, err := store.Authenticate(context.Background(), time.Now(), email, password)
	if err != nil {
		t.t.Fatal(err)
	}

	// generate token
	token, err := t.Auth.GenerateToken(claims)
	if err != nil {
		t.t.Fatal(err)
	}

	return token
}

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

// NewIntegration creates a database, seeds it and constructs an authenticator.
func NewIntegration(t *testing.T, dbc DbContainer) *Test {
	// init test environment
	logger, db, teardown := NewUnit(t, dbc)

	// create RSA keys to enable authentication
	keyId := uuid.NewString()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// build an authenticator using this private key for the key store
	keyMap := map[string]*rsa.PrivateKey{keyId: privateKey}
	authr, err := auth.New(keyId, keystore.NewMap(keyMap))
	if err != nil {
		t.Fatal(err)
	}

	return &Test{
		Db:       db,
		Logger:   logger,
		Auth:     authr,
		Teardown: teardown,
		t:        t,
	}
}

// StrPtr is a helper to return *string for a string.
func StrPtr(s string) *string {
	return &s
}

// IntPtr is a helper to return *int for an int.
func IntPtr(i int) *int {
	return &i
}
