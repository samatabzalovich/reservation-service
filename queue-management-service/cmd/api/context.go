package main

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

type User struct {
	ID        int64  `json:"id"`
	Activated bool   `json:"activated"`
	Type      string `json:"type"`
}

func (app *Config) contextSetUserId(r *http.Request, user *User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Config) contextGetUser(r *http.Request) (*User, error) {
	user, ok := r.Context().Value(userContextKey).(*User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
