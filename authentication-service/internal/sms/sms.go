package sms

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

func NewMessageService(url string, apiKey string) *MessageService {
	return &MessageService{
		Url:    url,
		ApiKey: apiKey,
	}
}

type MessageService struct {
	Url    string
	ApiKey string
}

type Sms struct {
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

func (service *MessageService) SendSmsCode(phone string, code int) error {
	sms, err := service.toJson(Sms{
		Recipient: phone,
		Text:      "Your code for ReserveHUB is " + strconv.Itoa(code),
	})
	if err != nil {
		return err
	}
	// Create a new POST request with the JSON data
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/Message/SendSmsMessage?apiKey=%s", service.Url, service.ApiKey), bytes.NewBuffer(sms))
	if err != nil {
		return err
	}

	// Add appropriate headers (e.g., Content-Type)
	req.Header.Set("Content-Type", "application/json")

	// Send the request using an http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and print the response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (service *MessageService) GenerateCode() int {
	// this function will generate a random code ranged from 1000 to 9999
	return rand.Intn(9000) + 1000
}
