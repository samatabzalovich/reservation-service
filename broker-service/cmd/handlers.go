package main

import (
	"errors"
	"net/http"
)

// RequestPayload describes the JSON that this service accepts as an HTTP Post request

type RequestPayload struct {
	Action string       `json:"action"`
	Auth   AuthPayload  `json:"auth,omitempty"`
	Reg    RegPayload   `json:"reg,omitempty"`
	Token  TokenPayload `json:"token,omitempty"`
	Sms    SmsPayload   `json:"sms,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.rpcErrorJson(w, err)
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
	default:
		app.rpcErrorJson(w, errors.New("unknown action"))
	}
}
