package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// An Account relates OAuth providers to users, to enable one user associating
// with multiple OAuth providers and accounts in the future.
//
// table: "accounts"
type Account struct {
	// ID is the primary key of the "accounts" table.
	ID string `json:"id" db:"id"`
	// UserID is a foreign key to the "users" table.
	UserID int32 `json:"userId" db:"user_id"`
	// AccessToken is the access token provided by the OAuth provider.
	AccessToken string `json:"accessToken" db:"access_token"`
}

// An AccountService is a connection to the database with methods
// for interacting with the "accounts" table.
type AccountService struct {
	db *sqlx.DB
}

// NewAccountService returns a pointer to AccountService.
func NewAccountService(db *sqlx.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

// GetAccountByID returns an Account, reports whether the account exists inside
// the database with the given id, and returns an error from the Get method.
//
// sql.ErrNoRows isn't treated as an error, but rather makes the function
// report the account doesn't exist.
func (as AccountService) GetAccountByID(ID int32) (Account, bool, error) {
	var account Account
	query := "select * from accounts where id = $1"
	err := as.db.Get(&account, query, ID)
	if err != sql.ErrNoRows {
		if err != nil {
			return Account{}, false, err
		}
		return account, true, nil
	}
	return Account{}, false, nil
}

// CreateAccount returns an Account and returns an error from the Get method.
//
// If successful, it inserts a new row into the "accounts" table with the given data.
func (as AccountService) CreateAccount(tx *sqlx.Tx, ID int32, userID int32, accessToken string) (Account, error) {
	var account Account
	query := "insert into accounts (id, user_id, access_token) values ($1, $2, $3) returning *"
	err := tx.Get(&account, query, ID, userID, accessToken)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}
