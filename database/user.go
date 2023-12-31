package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID     int32       `json:"id"`
	Name   pgtype.Text `json:"name"`
	Email  string      `json:"email"`
	Avatar pgtype.Text `json:"avatar"`
}

type UserService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		db: db,
	}
}
