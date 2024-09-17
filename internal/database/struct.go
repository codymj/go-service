package database

import (
	"database/sql"
	"time"
)

type Config struct {
	User            string
	Password        string
	Host            string
	Port            int
	Name            string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

type Connection struct {
	DB *sql.DB
}
