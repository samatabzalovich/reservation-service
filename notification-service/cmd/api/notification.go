package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func (app *Config) sendRequest(requestBody RequestBody) bool {

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error occurred in JSON marshal. Err: %s", err)
	}

	resp, err := http.Post(app.NotificationServiceHost, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error occurred sending request to the endpoint. Err: %s", err)
	}
	defer resp.Body.Close()

	return true
}
