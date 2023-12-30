package main

import (
	"log"
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/webdevfuel/projectmotor/db"
	"github.com/webdevfuel/projectmotor/template"
)

func main() {
	err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// NOTE ->> only for testing, remove after actual interactions with database
		message, err := db.GetMessage()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// <<- NOTE
		template.Dashboard.ExecuteWriter(pongo2.Context{
			"message": message,
		}, w)
	})
	http.ListenAndServe("localhost:3000", r)
}
