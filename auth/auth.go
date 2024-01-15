package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type UserKey struct{}

func SetUserSession(w http.ResponseWriter, r *http.Request, userID int32, session *sessions.Session) error {
	// Clear previous state and code from session values
	delete(session.Values, "state")
	delete(session.Values, "code")
	// Store new user on session values
	session.Values["userID"] = userID
	return session.Save(r, w)
}
