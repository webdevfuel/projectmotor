package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/flosch/pongo2/v6"
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

func generateCSRFToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
