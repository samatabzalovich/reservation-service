package main

import (
	"errors"
	"net/http"
)

func (app *Config) HandleInstitutionSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "createInstitution":
		app.CreateInstitutionViaGRpc(w, requestPayload)
	case "updateInstitution":
		app.UpdateInstitutionViaGRpc(w, requestPayload)
	case "deleteInstitution":
		app.DeleteInstitutionViaGRpc(w, requestPayload)
	case "getInstitution":
		app.GetInstitutionViaGRpc(w, requestPayload)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) ListInstitutions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PageSize int
		Page int
		SearchText string
		Sort string
		CategoryId int
	}
	qs := r.URL.Query()

	input.SearchText = app.readString(qs, "search_query", "")
	page ,err := app.readInt(qs, "page", 1)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	input.Page = page
	input.PageSize, err = app.readInt(qs, "page_size", 20,)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	input.Sort = app.readString(qs, "sort", "id")
	input.PageSize, err = app.readInt(qs, "category_id", 0,)
	if err != nil {
		app.errorJson(w, err)
		return	
	}
	filter := FilterPayload{
		PageSize: input.PageSize,
		Page: input.Page,
		SearchText: input.SearchText,
		Sort: input.Sort,
		CategoryId: int64(input.CategoryId),
	}

	if input.SearchText == "" {
		app.GetInstitutionsByCategoryViaGRpc(w, filter)
	}else {
		app.SearchInstitutionsViaGRpc(w, filter)
	}
}

