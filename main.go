package main

import (
	"log"
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/webdevfuel/projectmotor/database"
	"github.com/webdevfuel/projectmotor/handler"
	"github.com/webdevfuel/projectmotor/template"
)

var store = sessions.NewCookieStore([]byte("9eb88ac21007908b8a192a6b842ad097622882a84560ddf1ba0c27702836a9b4"))

func main() {
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	err = database.SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(handler.HandlerOptions{
		DB:    db,
		Store: store,
	})
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// NOTE ->> only for testing, remove after actual interactions with database
		message, err := database.GetMessage()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// <<- NOTE
		template.Dashboard.ExecuteWriter(pongo2.Context{
			"message": message,
		}, w)
	})
	r.Get("/login", h.Login)
	http.ListenAndServe("localhost:3000", r)
}
