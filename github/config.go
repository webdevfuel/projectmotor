package github

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Config is a oauth2 configuration with
// a client id, client secret, scopes and endpoint.
//
// It initializes the oauth2 package Config with provided values.
func newConfig() *oauth2.Config {
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	if clientId == "" {
		log.Fatal("environment variable GITHUB_CLIENT_ID must be present")
	}
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("environment variable GITHUB_CLIENT_SECRET must be present")
	}
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}

var Config = newConfig()
