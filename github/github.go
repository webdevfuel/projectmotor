package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Email struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

type User struct {
	ID int32 `json:"id"`
}

type Data struct {
	ID           int32
	PrimaryEmail string
}

type GitHubOAuth2 struct {
	accessToken string
}

func (g *GitHubOAuth2) fetchEmails() ([]Email, error) {
	var emails []Email
	request, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return []Email{}, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.accessToken))
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

func (g *GitHubOAuth2) fetchUser() (User, error) {
	var user User
	request, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return User{}, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.accessToken))
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

func (g *GitHubOAuth2) primaryEmail(emails []Email) (string, error) {
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

func (g *GitHubOAuth2) GetData() (Data, error) {
	emails, err := g.fetchEmails()
	if err != nil {
		return Data{}, err
	}
	primaryEmail, err := g.primaryEmail(emails)
	if err != nil {
		return Data{}, err
	}
	user, err := g.fetchUser()
	if err != nil {
		return Data{}, err
	}
	return Data{
		ID:           user.ID,
		PrimaryEmail: primaryEmail,
	}, nil
}

func New(accessToken string) *GitHubOAuth2 {
	return &GitHubOAuth2{
		accessToken: accessToken,
	}
}
