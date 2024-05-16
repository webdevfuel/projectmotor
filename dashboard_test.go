package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevfuel/projectmotor/test"
)

func TestDashboard(t *testing.T) {
	_, server := test.NewTestServer()
	defer server.Close()
	session, err := test.GetTestSession()
	if err != nil {
		t.Errorf("error getting test session %f", err)
	}

	t.Run("navigating to dashboard page redirects to and renders login page", func(t *testing.T) {
		res, _ := http.Get(fmt.Sprintf("%s/", server.URL))
		data, _ := io.ReadAll(res.Body)
		body := string(data)
		assert := assert.New(t)
		assert.Equal(200, res.StatusCode)
		assert.Contains(body, "Login with GitHub")
	})

	t.Run("navigating to dashboard page renders it", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/", server.URL), nil)
		cookie := fmt.Sprintf("_projectmotor_session=%s;", session)
		req.Header.Set("cookie", cookie)
		client := &http.Client{}
		res, _ := client.Do(req)
		data, _ := io.ReadAll(res.Body)
		body := string(data)
		assert := assert.New(t)
		assert.Equal(200, res.StatusCode)
		assert.Contains(body, "Welcome back,")
	})
}
