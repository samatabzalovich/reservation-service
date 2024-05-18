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
		case data.ErrInvalidInstId:
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
	instId, err := app.readIntParam(r, "instId")
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

func (app *Config) GetEmployeeScheduleAndService(w http.ResponseWriter, r *http.Request) {
	var (
		employeeID  int64
		serviceID   int64
		selectedDay time.Time
	)
	employeeID, err := app.readIntParam(r, "employee_id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	serviceID, err = app.readIntParam(r, "service_id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	selectedDay, err = app.readTimeParam("selected_day", r.URL.Query())
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	employee, err := app.Models.Employees.GetEmployeeScheduleAndService(employeeID, serviceID, selectedDay)
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

// Todo: check update method
func (app *Config) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	input, err := app.contextGetEmployee(r)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	err = app.Models.Employees.Update(input)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]any{"message": "employee updated"})
}

func (app *Config) UpdateEmployeeSchedule(w http.ResponseWriter, r *http.Request) {
	input,err := app.contextGetEmployee(r)

	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	if len(input.Schedule) == 0 || len(input.Schedule) > 7 {
		app.errorJson(w, data.ErrInvalidSchedule, http.StatusBadRequest)
		return
	}

	if input.ID <= 0 {
		app.errorJson(w, data.ErrInvalidEmployeeId, http.StatusBadRequest)
		return
	}

	
	schedules := make([]*data.EmployeeSchedule, 0, len(input.Schedule))

	for _, s := range input.Schedule {
		schedule, err := data.NewEmployeeSchedule(s.DayOfWeek, s.StartTime, s.EndTime, s.BreakStartTime, s.BreakEndTime)
		if err != nil {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}
		schedules = append(schedules, schedule)
	}
	input.Schedule =	schedules

	err = app.Models.Employees.UpdateSchedule(input)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"message": "employee schedule updated"})
}

func (app *Config) UpdateEmployeeServices(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmployeeID int64   `json:"employeeId"`
		Services   []int64 `json:"services"`
	}
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	if len(input.Services) == 0 {
		app.errorJson(w, data.ErrInvalidServices, http.StatusBadRequest)
		return
	}

	if input.EmployeeID <= 0 {
		app.errorJson(w, data.ErrInvalidEmployeeId, http.StatusBadRequest)
		return
	}

	employee, err := app.Models.Employees.GetById(input.EmployeeID)

	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorJson(w, data.ErrInvalidEmployeeId, http.StatusNotFound)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	owner, err := app.contextGetUser(r)

	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	inst, err := app.GetInstitution(employee.InstId)

	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	if inst.OwnerId != owner.ID {
		app.errorJson(w, data.ErrInvalidEmployeeORUserNotOwner, http.StatusForbidden)
		return
	}

	services := make([]*data.EmployeeServices, 0, len(input.Services))

	for _, s := range input.Services {
		service, err := data.NewEmployeeServices(s)
		if err != nil {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}
		services = append(services, service)
	}

	employee.Services = services

	err = app.Models.Employees.UpdateServices(employee)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorJson(w, data.ErrInvalidServiceId, http.StatusNotFound)
			return
		}
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]any{"message": "employee services updated"})
}

func (app *Config) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, err := app.readIntParam(r, "employee_id")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	//chekc if user is owner of institution
	inst, err := app.GetInstitutionForEmployee(employeeID)
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	user, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	if inst.OwnerId != user.ID {
		app.errorJson(w, ErrAuthentication, http.StatusForbidden)
		return
	}
	err = app.Models.Employees.Delete(employeeID)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]any{"message": "employee deleted"})
}
