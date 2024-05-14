package database

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sqlx.DB

const DATABASE_URL string = "postgres://emanuel@localhost:5432/projectmotor_dev?search_path=public&sslmode=disable"

func OpenDB() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("pgx", DATABASE_URL)
	if err != nil {
		return conn, err
	}
	DB = conn
	return conn, nil
}

func CloseDB() error {
	return DB.Close()
}
