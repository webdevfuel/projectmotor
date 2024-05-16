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
	_, s := test.NewTestServer()
	defer s.Close()

	t.Run("navigating to dashboard page redirects to and renders login page", func(t *testing.T) {
		res, _ := http.Get(fmt.Sprintf("%s/", s.URL))
		data, _ := io.ReadAll(res.Body)
		body := string(data)
		assert := assert.New(t)
		assert.Equal(200, res.StatusCode)
		assert.Contains(body, "Login with GitHub")
	})
}
