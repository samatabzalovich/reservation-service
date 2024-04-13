package main

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func (app *Config) GetCategories(w http.ResponseWriter, r *http.Request) {
	app.GetCategoriesViaGRpc(w)
}

func (app *Config) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	// get authentication header from request
	token, err := app.contextGetToken(r)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	tempMetadata := metadata.New(map[string]string{"authorization": token})
	ctx := metadata.NewOutgoingContext(context.Background(), tempMetadata)
	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, errors.New("invalid institution"))
		return
	}
	

	if requestPayload.Action == "createCategory" {
		app.CreateCategoryViaGRpc(w, ctx, requestPayload)
	}
}