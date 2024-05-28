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
		basePath = "/"
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

	// Add base path to routes
	base := mux.Route(basePath, func(r chi.Router) {
		r.Route("/employee", func(m chi.Router) {
			m.Use(app.requireAuthentication)
			m.Use(app.requireActivatedUser)
			m.Use(app.requireOwnerRole)
			m.Post("/create", app.CreateEmployee)
			m.Post("/qr-code", app.CreateQRCodeToken)
			m.Delete("/delete/{employee_id}", app.DeleteEmployee)
			m.Group(func(f chi.Router) {
				f.Use(app.requireInstOwnerForOwner)
				f.Put("/update", app.UpdateEmployee)
			})
			m.Group(func(f chi.Router) {
				f.Use(app.requireInstOwnerForOwnerToUpdateEmployeeSchedule)
				f.Put("/update-schedule", app.UpdateEmployeeSchedule) // TODO: check update method
			})
			m.Put("/update-services", app.UpdateEmployeeServices) //TODO: check update method
		})
	})

	base.Get("/schedule/{employee_id}/{service_id}", app.GetEmployeeScheduleAndService)

	base.Route("/ws/joinRegisterEmployeeRoom", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Get("/{token}", app.JoinRegisterEmployeeRoom)
	})

	base.Get("/institution-service/{id}", app.GetServiceForInstitution)
	base.Get("/institution-employee/{instId}", app.GetAllEmployeesForInstitution)

	base.Route("/service", func(r chi.Router) {
		r.Get("/{id}", app.GetService)
		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)
			r.Use(app.requireActivatedUser)
			r.Post("/create", app.CreateService) // TODO: add middleware to check if user is owner of institution
			r.Put("/update", app.UpdateService)
			r.Delete("/delete/{id}", app.DeleteService)
		})
	})

	base.Route("/institution", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Get("/inst-user-employee", app.GetInstitutionsUserEmployee)
	})

	mux.NotFound(app.NotFound)
	return mux
}
