package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// A UserKey is the representation of a key, for usage with context.
//
// An example of usage to set an user inside a context:
//
//	user := getUserById(0)
//	ctx := context.Background()
//	ctx = context.WithValue(ctx, auth.UserKey{}, user)
type UserKey struct{}

// SetUserSession returns an error when calling session.Save(r, w).
//
// It clears the previous values from the given session with keys "state" and "code".
//
// It sets the given userID on values with key "userID".
func SetUserSession(w http.ResponseWriter, r *http.Request, userID int32, session *sessions.Session) error {
	delete(session.Values, "state")
	delete(session.Values, "code")
	session.Values["userID"] = userID
	return session.Save(r, w)
}
