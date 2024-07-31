package github

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Config is a oauth2 configuration with
// a client id, client secret, scopes and endpoint.
//
// It initializes the oauth2 package Config with provided values.
var Config = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}
