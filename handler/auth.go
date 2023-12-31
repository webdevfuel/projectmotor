package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gorilla/sessions"
	"github.com/webdevfuel/projectmotor/template"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// NOTE ->> Replace values with environment variables in production
var config = &oauth2.Config{
	ClientID:     "98a6e6f797f9db728c6e",
	ClientSecret: "6391fd493bae19758df67538ae000a01172f6b9e",
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}

// <<- NOTE

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	template.Login.ExecuteWriter(pongo2.Context{}, w)
}

func (h Handler) OAuthGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateCSRFToken(16)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	session, err := h.GetSessionStore(r)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	url := config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h Handler) OAuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Get session store
	session, err := h.GetSessionStore(r)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	// Get state and code from url query (?state=foo&code=bar)
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	// Check if state matches between query and session
	if stateMatches(state, session) {
		// Exchange code for token
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
		log.Println("token: ", token.AccessToken)
		return
	}
	fail(w, err, http.StatusUnauthorized)
}

func generateCSRFToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func stateMatches(s string, session *sessions.Session) bool {
	state := session.Values["state"]
	if state == nil {
		return false
	}
	return s == state
}
