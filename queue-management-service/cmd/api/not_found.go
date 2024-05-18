package main

import (
	"errors"
	"log"
	"net/http"
)


func (app *Config) NotFound(w http.ResponseWriter, r *http.Request) {
	//print  the route that was not found to the console
	log.Printf("Route not found: %s", r.URL.Path)
	app.errorJson(w, errors.New("endpoint not found"), http.StatusNotFound)
}