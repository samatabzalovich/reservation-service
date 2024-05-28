package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:65326"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/", app.Broker)
	mux.Get("/health", app.HealthCheck)

	mux.Post("/handleAuth", app.HandleAuthSubmission)
	mux.Get("/institutions/owner/{ownerId}", app.GetInstitutionsForOwner)
	mux.Route("/handleInstitution", func(r chi.Router) {
		r.Use(app.Authenticate)
		r.Post("/", app.HandleInstitutionSubmission)

	})
	mux.Get("/institution/{instId}", app.GetInstitutionById)
	mux.Route("/category", func(r chi.Router) {
		r.Use(app.Authenticate)
		r.Post("/", app.CreateCategory)
		r.Delete("/{catId}", app.DeleteCategory)
		r.Put("/", app.UpdateCategory)
	})
	mux.Get("/listInstitutions", app.ListInstitutions)
	mux.Get("/getCategories", app.GetCategories)
	mux.Get("/getInstitutionCategories/{instId}", app.GetCategoriesForInstitution)

	return mux
}
