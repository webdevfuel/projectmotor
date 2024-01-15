package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Account struct {
	ID          string `json:"id" db:"id"`
	UserID      int32  `json:"userId" db:"user_id"`
	AccessToken string `json:"accessToken" db:"access_token"`
}

type AccountService struct {
	db *sqlx.DB
}

func NewAccountService(db *sqlx.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

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

func (as AccountService) CreateAccount(tx *sqlx.Tx, ID int32, userID int32, accessToken string) (Account, error) {
	var account Account
	query := "insert into accounts (id, user_id, access_token) values ($1, $2, $3) returning *"
	err := tx.Get(&account, query, ID, userID, accessToken)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}
