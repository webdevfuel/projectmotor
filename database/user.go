package database

import (
	"database/sql"

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

func (us UserService) GetUserByID(ID int32) (User, bool, error) {
	var user User
	query := "select * from users where id = $1"
	err := us.db.Get(&user, query, ID)
	if err != sql.ErrNoRows {
		if err != nil {
			return User{}, false, err
		}
		return user, true, nil
	}
	return User{}, false, nil
}

func (us UserService) CreateUser(tx *sqlx.Tx, email string) (User, error) {
	var user User
	query := "insert into users (email) values ($1) returning *"
	err := tx.Get(&user, query, email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
