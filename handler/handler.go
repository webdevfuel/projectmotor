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

// A Handler interacts with the database and cookie store.
type Handler struct {
	UserService    *database.UserService
	AccountService *database.AccountService
	ProjectService *database.ProjectService
	TaskService    *database.TaskService
	Store          *sessions.CookieStore
	DB             *sqlx.DB
}

// HandlerOptions is a representation of the options that should
// be passed to a handler when initialized.
type HandlerOptions struct {
	DB    *sqlx.DB
	Store *sessions.CookieStore
}

// NewHandler returns a new Handler.
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

// GetSessionStore returns a new Session and an error from the SessionStore Get method.
func (h *Handler) GetSessionStore(r *http.Request) (*sessions.Session, error) {
	return h.Store.Get(r, "_projectmotor_session")
}

// BeginTx returns a new Tx and an error from the DB BeginTxx method.
func (h *Handler) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return &sqlx.Tx{}, err
	}
	return tx, nil
}

// GetUserFromContext returns a User from the given context.
//
// The method never fails, and if the user doesn't exist within the request
// context, it returns a User struct with zero values.
//
// The method should only be used on handlers that run after a middleware
// that sets the user within the request context.
func (h *Handler) GetUserFromContext(ctx context.Context) database.User {
	user := ctx.Value(auth.UserKey{})
	if user, ok := user.(database.User); ok {
		return user
	}
	return database.User{}
}

// GetIDFromRequest returns the int32 value of a url param, extracted with
// the chi package.
//
// It uses the strconv.Atoi function internally to convert the url param
// from a string into an integer.
func (h *Handler) GetIDFromRequest(r *http.Request, key string) (int32, error) {
	id, err := strconv.Atoi(chi.URLParam(r, key))
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

// Error replies to the request with given HTTP code, and a status text given
// the HTTP code, and prints the error messago to the console.
func (h *Handler) Error(w http.ResponseWriter, err error, code int) {
	http.Error(w, http.StatusText(code), code)
	log.Println("error:", err)
}

// TriggerEvent joins the given slice of events with a comma-separated string
// and sets the result as a response header with the key "HX-Trigger". The method should be called
// in the context of an HTMX request.
func (h *Handler) TriggerEvent(w http.ResponseWriter, events ...string) {
	s := strings.Join(events, ", ")
	w.Header().Set("HX-Trigger", s)
}

// Reswap sets the given strategy as a response header with the key "HX-Reswap".
// The method should be called in the context of an HTMX request.
func (h *Handler) Reswap(w http.ResponseWriter, strategy string) {
	w.Header().Set("HX-Reswap", strategy)
}

// Redirect sets the given url as a response header with the key "HX-Redirect".
// The method should be called in the context of an HTMX request.
func (h *Handler) Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}

// ReplaceUrl sets the given url as a response header with the key "HX-Replace-Url".
// The method should be called in the context of an HTMX request.
func (h *Handler) ReplaceUrl(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Replace-Url", url)
}

// IsHTMXRequest reports whether the given request has the request header "Hx-Request" set.
func (h *Handler) IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("Hx-Request") != ""
}

// URLQuery is a representation of the standard library URL query, with a value
// and a way to check if it's empty.
type URLQuery struct {
	// Value given the query key.
	Value string
	// Reports whether the value is an empty string.
	IsEmpty bool
}

// GetURLQuery returns a new URLQuery given the request and key.
func (h *Handler) GetURLQuery(r *http.Request, key string) URLQuery {
	value := r.URL.Query().Get(key)

	return URLQuery{
		Value:   value,
		IsEmpty: value == "",
	}
}
