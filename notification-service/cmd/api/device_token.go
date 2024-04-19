package main

import (
	"net/http"
	"notification-service/internal/data"

	"github.com/go-chi/chi/v5"
)

func (app *Config) Insert(w http.ResponseWriter, r *http.Request) {
	var token data.DeviceToken
	err := app.readJSON(w, r, &token)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	result, err := app.Models.DeviceTokens.Insert(token)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "Device token inserted successfully", "id": result})
}

func (app *Config) GetByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	result, err := app.Models.DeviceTokens.GetByToken(token)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, result)
}

func (app *Config) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}
	result, err := app.Models.DeviceTokens.GetByUserID(userID)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, result)
}

func (app *Config) Update(w http.ResponseWriter, r *http.Request) {
	var token data.DeviceToken
	err := app.readJSON(w, r, &token)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	err = app.Models.DeviceTokens.Update(token)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]interface{}{"message": "Device token updated successfully"})
}

func (app *Config) DeleteByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.Models.DeviceTokens.DeleteByToken(token)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]interface{}{"message": "Device token deleted successfully"})
}

func (app *Config) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}
	err = app.Models.DeviceTokens.Delete(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]interface{}{"message": "Device token deleted successfully"})
}
