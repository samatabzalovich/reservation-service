package main

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func (app *Config) HandleInstitutionSubmission(w http.ResponseWriter, r *http.Request) {
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

	switch requestPayload.Action {
	case "createInstitution":
		app.CreateInstitutionViaGRpc(w, ctx, requestPayload)
	case "updateInstitution":
		app.UpdateInstitutionViaGRpc(w, ctx, requestPayload)
	case "deleteInstitution":
		app.DeleteInstitutionViaGRpc(w, ctx, requestPayload)
	case "getInstitution":
		app.GetInstitutionViaGRpc(w, requestPayload)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) ListInstitutions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PageSize   int
		Page       int
		SearchText string
		Sort       string
		Categories []int64
	}
	qs := r.URL.Query()

	input.SearchText = app.readString(qs, "search_query", "")
	page, err := app.readInt(qs, "page", 1)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	input.Page = page
	input.PageSize, err = app.readInt(qs, "page_size", 20)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	input.Sort = app.readString(qs, "sort", "id")
	input.Categories, err = app.readListOfIntValues(qs, "category_id", []int64{})
	if err != nil {
		app.errorJson(w, err)
		return
	}
	filter := FilterPayload{
		PageSize:   input.PageSize,
		Page:       input.Page,
		SearchText: input.SearchText,
		Sort:       input.Sort,
		Categories: input.Categories,
	}

	app.SearchInstitutionsViaGRpc(w, filter)

}

func (app *Config) GetInstitutionsForOwner(w http.ResponseWriter, r *http.Request) {
	ownerId, err := app.readIntParam(r, "ownerId")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.GetInstitutionsForOwnerViaGrpc(w, RequestPayload{Institution: InstPayload{OwnerId: ownerId}})
}


func (app *Config) GetInstitutionById(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	instId, err := app.readIntParam(r, "instId")
	if err != nil {
		app.errorJson(w, err)
		return
	}
	requestPayload.Institution.Id = instId

	app.GetInstitutionViaGRpc(w, requestPayload)
}
