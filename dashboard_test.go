package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevfuel/projectmotor/test"
)

func TestDashboard(t *testing.T) {
	handler, server := test.NewTestServer()
	defer server.Close()

	cookie, err := test.SetTestUserSession(server, 1)
	if err != nil {
		t.Errorf("error setting user test session %s", err)
		return
	}

	err = test.ResetAndSeedDB(handler.DB)
	if err != nil {
		t.Errorf("error resetting and seeding database %s", err)
		return
	}

	t.Run("navigating to dashboard page redirects to and renders login page", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(server.URL),
		)
		res := test.Do(req)
		body := test.Body(res)
		assert := assert.New(t)
		assert.Equal(200, res.StatusCode)
		assert.Contains(body, "Login with GitHub")
	})

	t.Run("navigating to dashboard page renders it", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(server.URL),
			test.WithAuthentication(test.Authenticated, cookie),
		)
		res := test.Do(req)
		body := test.Body(res)
		assert := assert.New(t)
		assert.Equal(200, res.StatusCode)
		assert.Contains(body, "Welcome back, hello@webdevfuel.com")
		assert.Contains(body, "Dashboard")
		assert.Contains(body, "Projects")
		assert.Contains(body, "Tasks")
		assert.Contains(body, "Log out")
	})
}
