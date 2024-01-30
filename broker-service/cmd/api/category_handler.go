package main

import "net/http"

func (app *Config) GetCategories(w http.ResponseWriter, r *http.Request) {
	app.GetCategoriesViaGRpc(w)
}