package main

import (
	"errors"
	"net/http"
)


func (app *Config) NotFound(w http.ResponseWriter, r *http.Request) {
	app.errorJson(w, errors.New("endpoint not found"), http.StatusNotFound)
}