package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/webdevfuel/projectmotor/database"
	"github.com/webdevfuel/projectmotor/handler"
	"github.com/webdevfuel/projectmotor/router"
)

func getCookieSessionKey() string {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("environment variable SESSION_KEY must be set")
		return ""
	}
	return sessionKey
}

var store = sessions.NewCookieStore([]byte(getCookieSessionKey()))

func main() {
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	h := handler.NewHandler(handler.HandlerOptions{
		DB:    db,
		Store: store,
	})
	r := router.NewRouter(h)
	http.ListenAndServe("localhost:3000", r)
}
