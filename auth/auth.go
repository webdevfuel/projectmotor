package auth

import (
	"crypto/rand"
	"encoding/base64"
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
// It sets the given token on values with key "token".
func SetUserSession(
	w http.ResponseWriter,
	r *http.Request,
	token string,
	session *sessions.Session,
) error {
	delete(session.Values, "state")
	delete(session.Values, "code")
	session.Values["token"] = token
	return session.Save(r, w)
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(b)
	return s, nil
}
