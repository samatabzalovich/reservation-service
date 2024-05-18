package main

import (
	"net/http"
	"strings"
)

func (app *Config) requireActivatedUser(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.contextGetUser(r)
		if err != nil {
			app.errorJson(w, err, http.StatusUnauthorized)
			return
		}
		if !user.Activated {
			app.errorJson(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Config) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			app.errorJson(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.errorJson(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		user, err := app.AuthenticateViaGrpc(token)
		if err != nil {
			app.rpcErrorJson(w, err)
			return
		}

		r = app.contextSetUserId(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *Config) requireAtLeastOneAppointmentOrQueue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.contextGetUser(r)
		if err != nil {
			app.errorJson(w, err, http.StatusUnauthorized)
			return
		}

		instId := app.readString(	r.URL.Query(), "instId", "")
		if instId == "" {
			app.errorJson(w, ErrNoInstId, http.StatusBadRequest)
			return
		}

		employeeId := app.readString(r.URL.Query(), "employeeId", "")
		

		appointments, err := app.GetAppointmentsForClientInInstitution(user.ID, instId, employeeId)

		if err != nil {
			app.errorJson(w, err, http.StatusInternalServerError)
			return
		}

		if appointments == 0 {
			queueCount, err := app.GetQueueForClientInInstitution(user.ID, instId, employeeId)
			if err != nil {
				app.errorJson(w, err, http.StatusInternalServerError)
				return
			}

			if queueCount == 0 {
				app.errorJson(w, ErrNoProvidedServiceInThatInstitution, http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
