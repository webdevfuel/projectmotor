package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Email is a representation of data returned from the GitHub API when
// making a GET request to "https://api.github.com/user/emails".
type Email struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// User is a representation of data returned from the GitHub API when
// making a GET request to "https://api.github.com/user".
type User struct {
	ID int32 `json:"id"`
}

// Data is a representation of a user's GitHub account.
type Data struct {
	// ID is the id of the user's GitHub account.
	ID int32
	// PrimaryEmail is the primary email address of the user's GitHub account.
	PrimaryEmail string
}

// GitHubOAuth2 is a wrapper around a GitHub OAuth access token.
//
// It's used to attach methods to it to interact with the GitHub API.
type GitHubOAuth2 struct {
	accessToken string
}

// GetData returns a new Data and the first error encountered when fetching
// the GitHub API.
//
// It fetches the GitHub API to know what the primary email and id of
// the user is given the access token.
func (g *GitHubOAuth2) GetData() (Data, error) {
	emails, err := fetchEmails(g.accessToken)
	if err != nil {
		return Data{}, err
	}
	primaryEmail, err := primaryEmail(emails)
	if err != nil {
		return Data{}, err
	}
	user, err := fetchUser(g.accessToken)
	if err != nil {
		return Data{}, err
	}
	return Data{
		ID:           user.ID,
		PrimaryEmail: primaryEmail,
	}, nil
}

// NewGitHubOAuth2 returns a new GitHubOAuth2 with the given access token.
func NewGitHubOAuth2(accessToken string) *GitHubOAuth2 {
	return &GitHubOAuth2{
		accessToken: accessToken,
	}
}

func fetchEmails(accessToken string) ([]Email, error) {
	var emails []Email
	request, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return []Email{}, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return []Email{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&emails)
	if err != nil {
		return []Email{}, err
	}
	return emails, nil

}

func fetchUser(accessToken string) (User, error) {
	var user User
	request, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return User{}, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return User{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func primaryEmail(emails []Email) (string, error) {
	var primary string
	for _, email := range emails {
		if email.Primary {
			primary = email.Email
		}
	}
	if primary == "" {
		return "", errors.New("no primary email found")
	}
	return primary, nil
}
