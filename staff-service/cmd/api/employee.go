package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"staff-service/internal/data"
	"time"
)

func (app *Config) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var input data.Employee
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(byteData, &input)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	err = app.Models.Employees.Insert(&input)
	if err != nil {
		log.Println("Error inserting employee: ")
		switch err {
		case data.ErrInvalidServices:
			app.errorJson(w, err, http.StatusBadRequest)
		default:
			app.errorJson(w, err, http.StatusInternalServerError)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, map[string]any{"message": "employee created", "id": input.ID})
}

func (app *Config) CreateQRCodeToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InstId int64 `json:"instId"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	owner, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}
	token, err := app.SetEmployeeRegTokenViaGrpc(owner.ID, input.InstId)

	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, map[string]string{"qr_token": token})
}

func (app *Config) GetAllEmployeesForInstitution(w http.ResponseWriter, r *http.Request) {
	instId,err := app.readIntParam(r, "instId")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	employees, err := app.Models.Employees.GetAllForInst(instId)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]any{"employees": employees})
}


func (app *Config) GetEmployeeScheduleAndService (w http.ResponseWriter, r *http.Request) {
	var (
		employeeID int64 
		serviceID int64 
		selectedDay time.Time 
	)
	employeeID, err := app.readIntParam(r, "employee_id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	serviceID,err = app.readIntParam(r, "service_id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	selectedDay,err = app.readTimeParam( "selected_day", r.URL.Query())
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	log.Println("employeeID: ", employeeID)
	log.Println("serviceID: ", serviceID)
	log.Println("selectedDay: ", selectedDay)
	employee, err := app.Models.Employees.GetEmployeeScheduleAndService( employeeID, serviceID, selectedDay)
	if err != nil {
		if err == data.ErrRecordNotFound {	
			app.errorJson(w, err, http.StatusNotFound)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, employee)
}

