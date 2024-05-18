package main

import (
	"errors"
	"net/http"
)

// RequestPayload describes the JSON that this service accepts as an HTTP Post request

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleAuthSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.AuthViaGRpc(w, requestPayload)
	case "token":
		app.CreateTokenViaGRpc(w, requestPayload)
	case "reg":
		app.RegViaGRpc(w, requestPayload)
	case "verifySms":
		app.VerifySmsViaGRpc(w, requestPayload)
	case "deleteAccount":
		app.DeleteAccountViaGRpc(w, requestPayload)
	case "sendSms":
		app.SendSmsViaGRpc(w, requestPayload)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}
