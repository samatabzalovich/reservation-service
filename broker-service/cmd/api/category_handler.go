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

func (app *Config) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// get authentication header from request
	token, err := app.contextGetToken(r)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	tempMetadata := metadata.New(map[string]string{"authorization": token})
	ctx := metadata.NewOutgoingContext(context.Background(), tempMetadata)
	categoryId, err := app.readIntParam(r, "catId")

	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	app.DeleteCategoryViaGRpc(w, ctx, categoryId)
}

func (app *Config) UpdateCategory(w http.ResponseWriter, r *http.Request) {
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

	if requestPayload.Action == "updateCategory" {
		app.UpdateCategoryViaGRpc(w, ctx, requestPayload)
	}
}

func (app *Config) GetCategoriesForInstitution(w http.ResponseWriter, r *http.Request) {
	// get institution id from request
	instId, err := app.readIntParam(r, "instId")
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	app.GetCategoriesForInstitutionViaGRpc(w, instId)
}
