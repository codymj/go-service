package database

import (
	"database/sql"
	"fmt"
)

func Configure(cfg *Config) (*Connection, error) {
	connectionStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return &Connection{DB: db}, nil
}
