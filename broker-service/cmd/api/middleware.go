package main

import (
	"errors"
	"net/http"
	"strings"
)

func (app *Config) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			app.errorJson(w, errors.New("authorization header required"))
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.errorJson(w, errors.New("authorization header format must be Bearer {token}"))
			return
		}

		token := headerParts[1]

		if len(token) != 26 {
			app.errorJson(w, errors.New("token should be 26 characters long"))
			return
		}

		r = app.contextSetToken(r, token)

		next.ServeHTTP(w, r)
	})
}
