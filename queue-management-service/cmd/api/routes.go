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
    mux.Get("/queue-number/{clientId}", app.GetQueueNumberForClientInInstitution)
    // Add base path to routes
    mux.Route(basePath + "/queue", func(r chi.Router) {
        r.Use(app.requireAuthentication)
        r.Use(app.requireActivatedUser)
        r.Post("/call-next", app.CallNextClient)
        r.Get("/join/{serviceId}", app.JoinQueueForServiceRoom)
        r.Get("/amount/{serviceId}", app.JoinQueueForPeopleAmountRoom)
        r.Get("/get-all-for-inst/{instId}", app.GetAllForInstitution)
        r.Put("/update-status", app.UpdateQueueStatus)
        r.Delete("/delete-all-for-inst/{instId}", app.DeleteAllForInst)
        r.Delete("/delete-by-id/{id}", app.DeleteQueue)
    })

    mux.NotFound(app.NotFound)
    return mux
}