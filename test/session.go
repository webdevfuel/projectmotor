package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
)

// SetTestUserSession returns a cookie with the format `_projectmotor_session=%s;`
// where %s is the session string with data about a user.
//
// It uses the test session.CookieStore and stores the given userId with key "userID".
func SetTestUserSession(userId int) (string, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://localhost:3000", nil)
	session, err := store.Get(r, "_projectmotor_session")
	if err != nil {
		return "", err
	}
	session.Values["userID"] = userId
	err = session.Save(r, w)
	if err != nil {
		return "", err
	}
	setCookie := w.Result().Header.Get("set-cookie")
	if setCookie == "" {
		return "", errors.New("set cookie cannot be an empty string")
	}
	parts := strings.Split(setCookie, " ")
	if len(parts) == 0 {
		return "", errors.New("set cookie has zero parts when split at empty space")
	}
	return parts[0], nil
}
