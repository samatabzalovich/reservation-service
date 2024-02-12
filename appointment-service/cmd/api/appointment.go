package main

import (
	"appointment-service/internal/data"
	"net/http"
	"time"
)

func (app *Config) GetAppointmentById(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointment, err := app.Models.Appointments.GetById(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map [string]any{"appointment": appointment})
}

func (app *Config) GetAppointmentsForInstitution(w http.ResponseWriter, r *http.Request) {
	institutionId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}
	
	appointments, err := app.Models.Appointments.GetAllForInst(institutionId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"appointments": appointments})
}

func (app *Config) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ClientID    int       `json:"client_id"`
	InstId      int       `json:"inst_id"`
	EmployeeID  int       `json:"employee_id"`
	ServiceID   int       `json:"service_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointment, err := data.NewAppointment(
		input.ClientID,
		input.InstId,
		input.EmployeeID,
		input.ServiceID,
		input.StartTime,
		input.EndTime,
		false,
	)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	err = app.Models.Appointments.Insert(appointment)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, map[string]any{"id": appointment.ID})
}


func (app *Config) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment data.Appointment
	err := app.readJSON(w, r, &appointment)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	err = app.Models.Appointments.Update(&appointment)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, appointment)
}

func (app *Config) GetAppointmentsForEmployee(w http.ResponseWriter, r *http.Request) {
	employeeId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointments, err := app.Models.Appointments.GetAllForEmployee(employeeId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map [string]any{"appointments": appointments})
}


func (app *Config) GetAppointmentsForClient(w http.ResponseWriter, r *http.Request) {
	clientId, err := app.readIntParam(r, "id")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	appointments, err := app.Models.Appointments.GetAllForClient(clientId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map [string]any{"appointments": appointments})
}