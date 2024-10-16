package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/webdevfuel/projectmotor/auth"
	"github.com/webdevfuel/projectmotor/handler"
)

// NewRouter returns a new chi.Mux router with all of the default middleware.
func NewRouter(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/login", h.Login)
	r.Get("/oauth/github/login", h.OAuthGitHubLogin)
	r.Get("/oauth/github/callback", h.OAuthGitHubCallback)
	r.Group(protectedRouter(h))
	return r
}

// Router with user ensured
//
// Add routes here where user has to be logged in
func protectedRouter(h *handler.Handler) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(protectedCtx(h))
		r.Get("/projects", h.GetProjects)
		r.Post("/projects", h.CreateProject)
		r.Get("/projects/new", h.NewProject)
		r.Get("/projects/{id}/edit", h.EditProject)
		r.Patch("/projects/{id}/toggle", h.ToggleProjectPublished)
		r.Patch("/projects/{id}", h.UpdateProject)
		r.Delete("/projects/{id}", h.DeleteProject)
		r.Get("/projects/{id}/share", h.ShareProject)
		r.Post("/projects/{id}/share", handler.ErrorWrapper(h.ShareProjectByEmail))
		r.Delete("/projects/{projectId}/share/{userId}", handler.ErrorWrapper(h.RevokeProjectById))
		r.Get("/tasks/new", h.NewTask)
		r.Post("/tasks", h.CreateTask)
		r.Get("/tasks", h.GetTasks)
		r.Get("/tasks/{id}/edit", h.EditTask)
		r.Patch("/tasks/{id}", h.UpdateTask)
		r.Get("/tasks/{id}", h.GetTask)
		r.Get("/", h.Dashboard)
	}
}

// Redirect to public auth route
//
// Use this when session user doesn't exist
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Protected context
//
// Middleware checks if user exists within current session
func protectedCtx(h *handler.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session, err := h.GetSessionStore(r)
			// redirect in case of error fetching session from store
			if err != nil {
				redirectToLogin(w, r)
				return
			}
			token := session.Values["token"]
			// redirect in case of missig token
			if token == nil {
				redirectToLogin(w, r)
				return
			}
			ctx := r.Context()
			// ensure token type is string
			tokenStr, ok := token.(string)
			if !ok {
				redirectToLogin(w, r)
				return
			}
			// check if user with token exists in DB
			user, err := h.UserService.GetUserBySessionToken(tokenStr)
			if err != nil {
				redirectToLogin(w, r)
				return
			}
			ctx = context.WithValue(ctx, auth.UserKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
