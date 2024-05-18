package sms

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

func NewMessageService(url string) *MessageService {
	return &MessageService{
		Url:    url,
	}
}

type MessageService struct {
	Url    string
}

type Sms struct {
	ChatID  string `json:"chatId"`
	Message string `json:"message"`
}

func (service *MessageService) SendSmsCodeAPI(phone string, code int) error {
	// remove + from phone number
	cleanedPhone := strings.ReplaceAll(phone, "+", "")
	sms, err := service.toJson(Sms{
		ChatID:  cleanedPhone + "@c.us",
		Message: "Your code for ReserveHUB is " + fmt.Sprintf("*%d*", code),
	})
	if err != nil {
		return err
	}
	// Create a new POST request with the JSON data
	req, err := http.NewRequest("POST", service.Url, bytes.NewBuffer(sms))
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
