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
			app.errorJson(w, ErrNotActivated, http.StatusForbidden)
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