package database

import "github.com/jmoiron/sqlx"

type Account struct {
	UserID      int32  `json:"userId"`
	AccountID   string `json:"accountId"`
	AccessToken string `json:"accessToken"`
}

type AccountService struct {
	db *sqlx.DB
}

func NewAccountService(db *sqlx.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}
