package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "/analytics-service"
	}

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
	base := mux.Route(basePath, func(u chi.Router) {
		u.Route("/comment", func(r chi.Router) {
			r.Use(app.requireAuthentication)
			r.Use(app.requireActivatedUser)
			r.Post("/", app.LeaveComment)
			r.Get("/institution/{id}", app.GetCommentsForInstitution)
			r.Get("/{id}", app.GetComment)
			r.Get("/user/{id}", app.GetCommentsForUser)
			r.Put("/", app.UpdateComment)
			r.Delete("/{id}", app.DeleteComment)
			r.Delete("/institution/{id}", app.DeleteCommentsForInstitution)
			r.Delete("/user/{id}", app.DeleteCommentsForUser)
		})
	})

	base.Route("/rating", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Post("/", app.LeaveFeedbackForAppointment)

		r.Put("/", app.UpdateFeedback)
		r.Get("/feedback-appointment/{id}", app.GetFeedbackForAppointment)
		r.Get("/employee/{id}", app.GetFeedbacksForEmployee)
		r.Get("/client/{id}", app.GetFeedbacksForClient)
		r.Get("/institution/{id}", app.GetFeedbacksForInstitution)

		r.Delete("/{id}", app.DeleteFeedback)
	})

	base.Route("/rating-analytics", func(r chi.Router) {
		r.Use(app.requireAuthentication)
		r.Use(app.requireActivatedUser)
		r.Get("/employee/{id}", app.GetAverageRatingForEmployee)
		r.Get("/client/{id}", app.GetAverageRatingForClient)
		r.Get("/institution/{id}", app.GetAverageRatingForInstitution)
		r.Get("/service/{id}", app.GetAverageRatingForService)
		r.Get("/employee-service/{employee_id}/{service_id}", app.GetAverageRatingForEmployeeService)
		r.Get("/client-service/{client_id}/{service_id}", app.GetAverageRatingForClientService)
		r.Get("/client-institution/{client_id}/{institution_id}", app.GetAverageRatingForClientInstitution)
		r.Get("/client-employee/{client_id}/{employee_id}", app.GetAverageRatingForClientEmployee)
	})

	mux.NotFound(app.NotFound)
	return mux
}
