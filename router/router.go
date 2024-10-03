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
		r.Get("/projects/{id}/share", h.ShareProject)
		r.Post("/projects/{id}/share", handler.FragmentWrapper(h.ShareProjectByEmail))
		r.Delete("/projects/{projectId}/share/{userId}", handler.FragmentWrapper(h.RevokeProjectById))
		r.Patch("/projects/{id}/toggle", h.ToggleProjectPublished)
		r.Patch("/projects/{id}", h.UpdateProject)
		r.Delete("/projects/{id}", h.DeleteProject)
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
			// redirect in case of error
			if err != nil {
				redirectToLogin(w, r)
				return
			}
			userID := session.Values["userID"]
			// redirect in case of missing user
			if userID == nil {
				redirectToLogin(w, r)
				return
			}
			ctx := r.Context()
			// check if userID type is int32
			if userID, ok := userID.(int32); ok {
				// check if user exists in DB
				user, exists, err := h.UserService.GetUserByID(userID)
				if err != nil {
					redirectToLogin(w, r)
					return
				}
				if !exists {
					redirectToLogin(w, r)
					return
				}
				ctx = context.WithValue(ctx, auth.UserKey{}, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
