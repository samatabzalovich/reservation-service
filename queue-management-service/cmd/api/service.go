package main

import (
	"fmt"
	"time"
)

type Service struct {
	ID          int64         `json:"id"`
	InstId      int64         `json:"inst_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       int           `json:"price"`
	Duration    time.Duration `json:"duration"`
	PhotoUrl    string        `json:"photo_url"`
	CreatedAt   string        `json:"created_at"`
	ServiceType string        `json:"serviceType"`
	UpdatedAt   string        `json:"updated_at"`
}

func (app *Config) GetServiceById(id int64) (*Service, error) {
	var service Service
	if (app.staffServiceHost == "") {
		app.staffServiceHost = "localhost:8082/staff-services-host/"
	}
	endpoint := fmt.Sprintf("http://%s%s%d", app.staffServiceHost, "service/", id)
	err := app.sendGetRequest(endpoint, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}
// aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin  982796488014.dkr.ecr.us-east-1.amazonaws.com
