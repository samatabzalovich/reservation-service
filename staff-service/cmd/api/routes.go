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
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Heartbeat("/ping"))

	// mux.Get("/employee/{id}", app.GetEmployeeById)
	// mux.Get("/institution-employees/{id}", app.GetEmployeesForInstitution)
	mux.Route("/employee", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Post("/create", app.CreateEmployee)
		r.Post("/qr-code", app.CreateQRCodeToken)
		r.Get("/{instId}", app.GetAllForInstitution)

		// r.Put("/update", app.)
		// r.Delete("/delete", app.)
		// r.Put("/update-services", app.)
		// r.Put("/update-schedule", app.)
	})

	mux.Route("/ws/joinRegisterEmployeeRoom", func(r chi.Router) {

		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Get("/{token}", app.JoinRegisterEmployeeRoom)
	})

	//mux.Get("/service/{id}", app.GetServiceById)
	//mux.Get("/institution-services/{id}", app.GetServicesForInstitution)
	mux.Route("/service", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Post("/create", app.CreateService) //TODO: add middleware to check if user is owner of institution
		//r.Put("/update", app.)
		//r.Delete("/delete/{id}", app.)
	})
	//mux.NotFound(app.NotFound)
	return mux
}
