package main

import (
	"log"
	"net/http"
)

func (app *Config) GetInstitutionsUserEmployee(w http.ResponseWriter, r *http.Request) {
	employeeUser, err := app.contextGetUser(r)
	if err != nil {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}
	institutions, err := app.GetInstitutionsForUserEmployee(employeeUser.ID)
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	log.Println("Institutions for user employee: ", institutions)
	app.writeJSON(w, http.StatusOK, institutions)
}