package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/codymj/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidId   = errors.New("ID is not in its proper form")
	ErrAuthFailure = errors.New("authentication failed")
	ErrForbidden   = errors.New("attempted action is not allowed")
)

// Config contains required properties to use the database
type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTls   bool
}

// Open a database connection based on the configuration.
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTls {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslMode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	// ping db
	var pingErr error
	for attempts := 1; ; attempts++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// make sure we did not timeout or get cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// run a simple query to determine connectivity
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// NamedExecContext is a helper to execute a CUD operation with logging and tracing.
func NamedExecContext(ctx context.Context, logger *zerolog.Logger, db *sqlx.DB, query string, data any) error {
	// log query
	q := queryString(query, data)
	logger.Info().Timestamp().
		Str("traceId", web.GetTraceId(ctx)).
		Str("query", q).
		Msg("database.NamedExecContext")

	// run query
	if _, err := db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	return nil
}

// NamedQuerySlice is a helper to execute queries that return a collection of data.
func NamedQuerySlice(ctx context.Context, logger *zerolog.Logger, db *sqlx.DB, query string, data, dest any) error {
	// log query
	q := queryString(query, data)
	logger.Info().Timestamp().
		Str("traceId", web.GetTraceId(ctx)).
		Str("query", q).
		Msg("database.NamedQuerySlice")

	// ensure destination is a pointer to a slice
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	// run query
	rows, err := db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}

	// build result
	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err = rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}

// NamedQueryStruct is a helper for executing queries that return a single value.
func NamedQueryStruct(ctx context.Context, logger *zerolog.Logger, db *sqlx.DB, query string, data, dest any) error {
	// log query
	q := queryString(query, data)
	logger.Info().Timestamp().
		Str("traceId", web.GetTraceId(ctx)).
		Str("query", q).
		Msg("database.NamedQueryStruct")

	// run query
	rows, err := db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return ErrNotFound
	}

	// build result
	if err = rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// queryString provides a pretty print version of the query and parameters.
func queryString(query string, args ...any) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.Trim(query, " ")
}
