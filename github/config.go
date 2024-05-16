package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Config is a oauth2 configuration with
// a client id, client secret, scopes and endpoint.
//
// It initializes the oauth2 package Config with provided values.
var Config = &oauth2.Config{
	ClientID:     "98a6e6f797f9db728c6e",
	ClientSecret: "6391fd493bae19758df67538ae000a01172f6b9e",
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}
