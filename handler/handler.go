package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) GetSessionStore(r *http.Request) (*sessions.Session, error) {
	return h.Store.Get(r, "_projectmotor_session")
}

func (h *Handler) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return &sqlx.Tx{}, err
	}
	return tx, nil
}

func (h *Handler) GetUserFromContext(ctx context.Context) database.User {
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

func (h *Handler) Error(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println("error:", err)
}

func (h *Handler) TriggerEvent(w http.ResponseWriter, events ...string) {
	s := strings.Join(events, ", ")
	w.Header().Set("HX-Trigger", s)
}

func (h *Handler) Reswap(w http.ResponseWriter, strategy string) {
	w.Header().Set("HX-Reswap", strategy)
}

func (h *Handler) Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}

func (h *Handler) ReplaceUrl(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Replace-Url", url)
}

func (h *Handler) IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("Hx-Request") != ""
}

type URLQuery struct {
	Value   string
	IsEmpty bool
}

func (h *Handler) GetURLQuery(r *http.Request, key string) URLQuery {
	value := r.URL.Query().Get(key)

	return URLQuery{
		Value:   value,
		IsEmpty: value == "",
	}
}
