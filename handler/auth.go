package handler

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/auth"
	"github.com/webdevfuel/projectmotor/database"
	"github.com/webdevfuel/projectmotor/github"
	"github.com/webdevfuel/projectmotor/template"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	component := template.Login()
	component.Render(r.Context(), w)
}

func (h *Handler) OAuthGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateCSRFToken(16)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	session, err := h.GetSessionStore(r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	url := github.Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) OAuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Get session store
	session, err := h.GetSessionStore(r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Get state and code from url query (?state=foo&code=bar)
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	// Ensure state matches between query and session
	if !stateMatches(state, session) {
		h.Error(w, errors.New("session and query state mismatch"), http.StatusBadRequest)
		return
	}
	// Exchange code for token
	token, err := github.Config.Exchange(context.Background(), code)
	if err != nil {
		h.Error(w, err, http.StatusBadRequest)
		return
	}
	// Ensure token is valid
	if !token.Valid() {
		h.Error(w, err, http.StatusBadRequest)
		return
	}
	// Initialise github.GitHubOAuth2 instance
	gh := github.NewGitHubOAuth2(token.AccessToken)
	// Fetch data from GitHub's API
	data, err := gh.GetData()
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Check if user already exists
	exists, err := h.UserService.UserExistsByGitHubID(data.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Begin transaction
	tx, err := h.BeginTx(r.Context())
	defer tx.Rollback()
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Create or update user access token and email
	user, err := createOrUpdateUser(
		tx,
		exists,
		data.PrimaryEmail,
		token.AccessToken,
		data.ID,
		h.UserService,
	)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Generate a random token to use as session identifier
	sessionToken, err := auth.GenerateSessionToken()
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Create session
	err = h.SessionService.CreateToken(tx, user.ID, sessionToken)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// Set session on cookies with token
	err = auth.SetUserSession(w, r, sessionToken, session)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createOrUpdateUser(
	tx *sqlx.Tx,
	exists bool,
	primaryEmail string,
	accessToken string,
	id int32,
	userService *database.UserService,
) (database.User, error) {
	if exists {
		return userService.UpdateUser(tx, primaryEmail, accessToken, id)
	}
	return userService.CreateUser(tx, primaryEmail, accessToken, id)
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
	stateStr, ok := state.(string)
	if !ok {
		return false
	}
	return secureCompare(s, stateStr)
}

func secureCompare(given string, actual string) bool {
	givenLen := int32(len(given))
	actualLen := int32(len(actual))
	if subtle.ConstantTimeEq(givenLen, actualLen) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
	}
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sess, err := h.GetSessionStore(r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	tok, ok := sess.Values["token"].(string)
	if !ok {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	err = h.SessionService.DeleteToken(tok)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	delete(sess.Values, "token")
	err = sess.Save(r, w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.Redirect(w, "http://localhost:3000/login")
}
