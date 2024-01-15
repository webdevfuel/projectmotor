package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var Config = &oauth2.Config{
	ClientID:     "98a6e6f797f9db728c6e",
	ClientSecret: "6391fd493bae19758df67538ae000a01172f6b9e",
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}
