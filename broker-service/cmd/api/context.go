package main

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

const tokenContextKey = contextKey("token")

func (app *Config) contextSetToken(r *http.Request, token string) *http.Request {
	newCtx := context.WithValue(r.Context(), tokenContextKey, token)
	return r.WithContext(newCtx)
}

func (app *Config) contextGetToken(r *http.Request) (string, error) {
	user, ok := r.Context().Value(tokenContextKey).(string)
	if !ok {
		return "", errors.New("token not found in context")
	}
	return user, nil
}
