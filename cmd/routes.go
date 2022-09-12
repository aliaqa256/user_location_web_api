package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/", app.RootHandler)
	mux.Post("/user", app.CreateUserHandler)
	mux.Post("/user/info",app.AddUserInfo)
	mux.Get("/user/lastlocation/{id}",app.GetLastLocation)
	mux.Post("/user/pastlocations",app.GetPastLocations)


	return mux
}
