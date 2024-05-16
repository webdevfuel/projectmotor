package test

import (
	"net/http/httptest"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/handler"
	"github.com/webdevfuel/projectmotor/router"
)

var store *sessions.CookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// NewTestServer returns a new httptest.Server
//
// Usually initialized before a set of requests as part of a test.
func NewTestServer(db *sqlx.DB) *httptest.Server {
	h := handler.NewHandler(handler.HandlerOptions{
		DB:    db,
		Store: store,
	})
	r := router.NewRouter(h)
	return httptest.NewServer(r)
}
