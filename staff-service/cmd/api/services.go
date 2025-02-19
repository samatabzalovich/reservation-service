package main

import (
	"errors"
	"net/http"
	"staff-service/internal/data"
	"time"
)

func (app *Config) CreateService(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		InstID      int64  `json:"inst_id"`
		Price       int    `json:"price"`
		Description string `json:"description"`
		Duration    string `json:"duration"`
		PhotoUrl    string `json:"photo_url"`
		ServiceType string `json:"serviceType"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	duration, err := time.ParseDuration(input.Duration)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	service, err := data.NewService(input.InstID, input.Name, input.Description, input.Price, duration, input.PhotoUrl, input.ServiceType)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Service.Insert(service)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidInstId):
			app.errorJson(w, err, http.StatusBadRequest)
		case errors.Is(err, data.ErrInvalidServices):
			app.errorJson(w, err, http.StatusBadRequest)
		case errors.Is(err, data.ErrUserIsNotEmployee):
			app.errorJson(w, err, http.StatusBadRequest)
		default:
			app.errorJson(w, err, http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, map[string]int64{"id": service.ID})
}

func (app *Config) GetServiceForInstitution(w http.ResponseWriter, r *http.Request) {
	institutionId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	services, err := app.Models.Service.GetAllForInst(institutionId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, services)
}

func (app *Config) GetService(w http.ResponseWriter, r *http.Request) {
	serviceId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	service, err := app.Models.Service.GetById(serviceId)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorJson(w, err, http.StatusNotFound)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, service)
}

func (app *Config) UpdateService(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Price       int    `json:"price"`
		Description string `json:"description"`
		Duration    string `json:"duration"`
		PhotoUrl    string `json:"photoUrl"`
		ServiceType string `json:"serviceType"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	duration, err := time.ParseDuration(input.Duration)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	service := &data.Service{
		ID:          input.ID,
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		Duration:    duration,
		PhotoUrl:    input.PhotoUrl,
		ServiceType: input.ServiceType,
		InstId:      0,
	}

	err = app.Models.Service.Update(service)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorJson(w, err, http.StatusBadRequest)
		default:
			app.errorJson(w, err, http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Service updated"})
}

func (app *Config) DeleteService(w http.ResponseWriter, r *http.Request) {
	serviceId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Service.Delete(serviceId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Service deleted"})
}
