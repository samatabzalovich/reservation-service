package main

import (
	"net/http"
	"staff-service/internal/data"
)

func (app *Config) CreateService(w http.ResponseWriter, r *http.Request) {
	var input data.Service
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Service.Insert(&input)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, input.ID)
}
