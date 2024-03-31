package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/auth"
	"github.com/webdevfuel/projectmotor/database"
)

type Handler struct {
	UserService    *database.UserService
	AccountService *database.AccountService
	ProjectService *database.ProjectService
	TaskService    *database.TaskService
	Store          *sessions.CookieStore
	DB             *sqlx.DB
}

type HandlerOptions struct {
	DB    *sqlx.DB
	Store *sessions.CookieStore
}

func NewHandler(options HandlerOptions) *Handler {
	userService := database.NewUserService(options.DB)
	accountService := database.NewAccountService(options.DB)
	projectService := database.NewProjectService(options.DB)
	taskService := database.NewTaskService(options.DB)
	return &Handler{
		Store:          options.Store,
		DB:             options.DB,
		UserService:    userService,
		AccountService: accountService,
		ProjectService: projectService,
		TaskService:    taskService,
	}
}

func (h Handler) GetSessionStore(r *http.Request) (*sessions.Session, error) {
	return h.Store.Get(r, "_projectmotor_session")
}

func (h Handler) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return &sqlx.Tx{}, err
	}
	return tx, nil
}

func (h Handler) GetUserFromContext(ctx context.Context) database.User {
	user := ctx.Value(auth.UserKey{})
	if user, ok := user.(database.User); ok {
		return user
	}
	return database.User{}
}

func (h *Handler) GetIDFromRequest(r *http.Request, key string) (int32, error) {
	id, err := strconv.Atoi(chi.URLParam(r, key))
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

func fail(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println("error:", err)
}
