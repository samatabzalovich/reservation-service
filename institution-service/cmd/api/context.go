package main

import (
	"context"
	"errors"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Config) contextSetUser(ctx context.Context,user *AuthPayload) *context.Context {
	newCtx := context.WithValue(ctx, userContextKey, user)
	return &newCtx
}

func contextGetUser(ctx context.Context) (*AuthPayload, error){
	user, ok := ctx.Value(userContextKey).(*AuthPayload)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
