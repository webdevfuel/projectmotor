package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/database"
)

type Handler struct {
	userService *database.UserService
	store       *sessions.CookieStore
}

type HandlerOptions struct {
	DB    *sqlx.DB
	Store *sessions.CookieStore
}

func NewHandler(options HandlerOptions) *Handler {
	userService := database.NewUserService(options.DB)
	return &Handler{
		store:       options.Store,
		userService: userService,
	}
}

func (h Handler) GetSessionStore(r *http.Request) (*sessions.Session, error) {
	return h.store.Get(r, "_projectmotor_session")
}

func fail(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println("error:", err)
}
