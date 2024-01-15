package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/database"
)

type Handler struct {
	userService    *database.UserService
	accountService *database.AccountService
	store          *sessions.CookieStore
	db             *sqlx.DB
}

type HandlerOptions struct {
	DB    *sqlx.DB
	Store *sessions.CookieStore
}

func NewHandler(options HandlerOptions) *Handler {
	userService := database.NewUserService(options.DB)
	accountService := database.NewAccountService(options.DB)
	return &Handler{
		store:          options.Store,
		db:             options.DB,
		userService:    userService,
		accountService: accountService,
	}
}

func (h Handler) GetSessionStore(r *http.Request) (*sessions.Session, error) {
	return h.store.Get(r, "_projectmotor_session")
}

func (h Handler) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		return &sqlx.Tx{}, err
	}
	return tx, nil
}

func fail(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println("error:", err)
}
