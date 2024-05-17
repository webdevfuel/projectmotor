package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
)

func SetTestUserSession() (string, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://localhost:3000", nil)
	session, err := store.Get(r, "_projectmotor_session")
	if err != nil {
		return "", err
	}
	session.Values["userID"] = 1
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
