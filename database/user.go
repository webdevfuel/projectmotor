package database

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

// A User is where we keep the information about a person, with the unique
// point of data being the email address.
//
// The email address is not set by the user when signing up, but rather
// extracted from the OAuth provider if a user doesn't yet exist with
// the email address sent by the OAuth provider.
//
// table: "users"
type User struct {
	ID     int32       `json:"id"`
	Name   pgtype.Text `json:"name"`
	Email  string      `json:"email"`
	Avatar pgtype.Text `json:"avatar"`
}

// A UserService is a connection to the database with methods
// for interacting with the "users" table.
type UserService struct {
	db *sqlx.DB
}

// NewUserService returns a pointer to UserService.
func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetAccountByID returns a User, reports whether the user exists inside
// the database with the given id, and returns an error from the Get method.
//
// sql.ErrNoRows isn't treated as an error, but rather makes the function
// report the user doesn't exist.
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

// CreateUser returns a User and returns an error from the Get method.
//
// If successful, it inserts a new row into the "users" table with the given data.
func (us UserService) CreateUser(tx *sqlx.Tx, email string) (User, error) {
	var user User
	query := "insert into users (email) values ($1) returning *"
	err := tx.Get(&user, query, email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
