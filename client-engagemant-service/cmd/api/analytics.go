package main

import (
	"client-engagemant-service/internal/data"
	"net/http"
	"time"
)

func (app *Config) TotalAppointmentsOfInstitutionForGivenDateRange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InstitutionID int64     `json:"institutionId"`
		StartDate     time.Time `json:"startDate"`
		EndDate       time.Time `json:"endDate"`
	}
	var err error
	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	startDateString := app.readString(r.URL.Query(), "startDate", "")

	if startDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.StartDate, err = time.Parse(time.RFC3339, startDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	endDateString := app.readString(r.URL.Query(), "endDate", "")

	if endDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.EndDate, err = time.Parse(time.RFC3339, endDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	totalAppointments, err := app.Models.Analytics.TotalAppointmentsOfInstitutionForGivenDateRange(input.InstitutionID, input.StartDate, input.EndDate)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	type response struct {
		TotalAppointments int `json:"total_appointments"`
	}

	app.writeJSON(w, http.StatusOK, response{TotalAppointments: totalAppointments})
}

func (app *Config) WageOfEmployeeServiceForGivenDateRange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmployeeID int64     `json:"employeeId"`
		StartDate  time.Time `json:"startDate"`
		EndDate    time.Time `json:"endDate"`
	}
	var err error
	input.EmployeeID, err = app.readIntParam(r, "employeeId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	startDateString := app.readString(r.URL.Query(), "startDate", "")

	if startDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.StartDate, err = time.Parse(time.RFC3339, startDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	endDateString := app.readString(r.URL.Query(), "endDate", "")

	if endDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.EndDate, err = time.Parse(time.RFC3339, endDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	if input.EmployeeID < 1 {
		app.errorJson(w, data.ErrInvalidEmployeeId, http.StatusBadRequest)
		return
	}

	wage, err := app.Models.Analytics.WageOfEmployeeServiceForGivenDateRange(input.EmployeeID, input.StartDate, input.EndDate)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"wage": wage})
}

func (app *Config) TotalRevenueOfInstitutionForGivenDateRange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InstitutionID int64     `json:"institutionId"`
		StartDate     time.Time `json:"startDate"`
		EndDate       time.Time `json:"endDate"`
	}
	var err error
	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	startDateString := app.readString(r.URL.Query(), "startDate", "")

	if startDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.StartDate, err = time.Parse(time.RFC3339, startDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	endDateString := app.readString(r.URL.Query(), "endDate", "")

	if endDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.EndDate, err = time.Parse(time.RFC3339, endDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	totalRevenue, err := app.Models.Analytics.TotalRevenueOfInstitutionForGivenDateRange(input.InstitutionID, input.StartDate, input.EndDate)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"totalRevenue": totalRevenue})
}

func (app *Config) TotalAppointmentsPerEmployeeForGivenDateRange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		EmployeeID    int64     `json:"employeeId"`
		InstitutionID int64     `json:"institutionId"`
		StartDate     time.Time `json:"startDate"`
		EndDate       time.Time `json:"endDate"`
	}
	var err error
	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	input.EmployeeID, err = app.readIntParam(r, "employeeId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return

	}

	startDateString := app.readString(r.URL.Query(), "startDate", "")

	if startDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.StartDate, err = time.Parse(time.RFC3339, startDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	endDateString := app.readString(r.URL.Query(), "endDate", "")

	if endDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	input.EndDate, err = time.Parse(time.RFC3339, endDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	if input.EmployeeID < 1 {
		app.errorJson(w, data.ErrInvalidEmployeeId, http.StatusBadRequest)
		return
	}

	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	totalAppointments, err := app.Models.Analytics.TotalAppointmentsPerEmployeeForGivenDateRange(input.EmployeeID, input.InstitutionID, input.StartDate, input.EndDate)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"totalAppointments": totalAppointments})
}

func (app *Config) MostPopularServicesByAppointments(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ServiceID     int64 `json:"serviceId"`
		InstitutionID int64 `json:"institutionId"`
	}
	var err error
	input.ServiceID, err = app.readIntParam(r, "serviceId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return

	}

	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return

	}

	if input.ServiceID < 1 {
		app.errorJson(w, data.ErrInvalidServiceId, http.StatusBadRequest)
		return
	}

	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	mostPopularServices, err := app.Models.Analytics.MostPopularServicesByAppointments(input.ServiceID, input.InstitutionID)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"mostPopularServices": mostPopularServices})
}

func (app *Config) MostPopularAppointmentsBySelectedDate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InstitutionID int64  `json:"institutionId"`
		Date          string `json:"date"`
	}
	var err error

	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	allowedDates := []string{"year", "month", "day"}

	input.Date = app.readString(r.URL.Query(), "date", "year")

	if input.Date == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	if !contains(allowedDates, input.Date) {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}

	if input.Date == "day" {
		input.Date = "dow"
	}

	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	mostPopularAppointments, err := app.Models.Analytics.MostPopularAppointmentsBySelectedDate(input.InstitutionID, input.Date)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"mostPopularAppointments": mostPopularAppointments})
}

func (app *Config) GetCanceledAndTotalAppointmentsForGivenDateRange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InstitutionID int64     `json:"institutionId"`
		StartDate     time.Time `json:"startDate"`
		EndDate       time.Time `json:"endDate"`
	}

	var err error
	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return

	}

	startDateString := app.readString(r.URL.Query(), "startDate", "")

	if startDateString == "" {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	input.StartDate, err = time.Parse(time.RFC3339, startDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	endDateString := app.readString(r.URL.Query(), "endDate", "")

	if endDateString == "" {

		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return

	}

	input.EndDate, err = time.Parse(time.RFC3339, endDateString)
	if err != nil {
		app.errorJson(w, data.ErrInvalidDate, http.StatusBadRequest)
		return
	}
	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	appointments, err := app.Models.Analytics.GetCanceledAndTotalAppointmentsForGivenDateRange(input.InstitutionID, input.StartDate, input.EndDate)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"appointments": appointments})
}

func (app *Config) EmployeeWithHighestRating(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InstitutionID int64 `json:"institutionId"`
	}
	var err error
	input.InstitutionID, err = app.readIntParam(r, "institutionId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	if input.InstitutionID < 1 {
		app.errorJson(w, data.ErrInvalidInstitutionId, http.StatusBadRequest)
		return
	}

	employees, err := app.Models.Analytics.EmployeeWithHighestRating(input.InstitutionID)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{"employees": employees})
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
