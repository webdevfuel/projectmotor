package database

import (
	"errors"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenDB returns a pointer to sqlx.DB and the first encountered error
// when attemping to establish a connection to the database url.
//
// Environment variable DATABASE_URL must be set.
func OpenDB() (*sqlx.DB, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, errors.New("environment variable DATABASE_URL must be set")
	}
	conn, err := sqlx.Connect("pgx", databaseUrl)
	if err != nil {
		return conn, err
	}
	return conn, nil
}
