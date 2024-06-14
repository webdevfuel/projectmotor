package test

import (
	"log"
	"net/http/httptest"
	"os"

	"github.com/gorilla/sessions"
	"github.com/webdevfuel/projectmotor/database"
	"github.com/webdevfuel/projectmotor/handler"
	"github.com/webdevfuel/projectmotor/router"
)

var store *sessions.CookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// NewServer returns a new handler.Handler and httptest.Server
//
// Usually initialized before a set of requests as part of a test.
//
// We return the handler to aid with doing assertions on the database, since
// it's easier than creating abstractions just for testing.
func NewServer() (*handler.Handler, *httptest.Server) {
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(handler.HandlerOptions{
		DB:    db,
		Store: store,
	})
	r := router.NewRouter(h)
	return h, httptest.NewServer(r)
}
