package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "/"
	}
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/health", app.HealthCheck)
	mux.Route(basePath + "/notification", func(r chi.Router) {
		r.Post("/device-token", app.Insert)
		r.Get("/device-token/{token}", app.GetByToken)
		r.Get("/device-token/user/{id}", app.GetByUserID)
		r.Put("/device-token", app.Update)
		r.Delete("/device-token/{id}", app.Delete)
	})
	mux.NotFound(app.NotFound)
	return mux
}


func (app *Config) NotFound(w http.ResponseWriter, r *http.Request) {
	//print  the route that was not found to the console
	log.Printf("Route not found: %s", r.URL.Path)
	app.errorJson(w, errors.New("endpoint not found"), http.StatusNotFound)
}


// TODO: get rid off unnecessary logs