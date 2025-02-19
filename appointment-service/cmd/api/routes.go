package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = ""
	}
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
	mux.Get("/health", app.HealthCheck)
	mux.Get(basePath + "/available-time-slots/{employee_id}", app.GetAvailableTimeSlots)
	mux.Get(basePath + "/completed-appointments-number/{clientId}", app.GetNumberOfCompletedAppointmentsForUser)
	mux.Route(basePath + "/appointment", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Get("/{id}", app.GetAppointmentById)
		r.Get("/institution-appointments/{id}", app.GetAppointmentsForInstitution)
		r.Get("/employee-appointments/{id}", app.GetAppointmentsForEmployee)
		r.Get("/client-appointments/{id}", app.GetAppointmentsForClient)
		r.Post("/create", app.CreateAppointment)
		r.Put("/update", app.UpdateAppointment)
	})
	mux.NotFound(app.NotFound)
	return mux
}
