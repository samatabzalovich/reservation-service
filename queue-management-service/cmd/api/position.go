package main

import (
	"errors"
	"queue-managemant-service/internal/data"
)

func (app *Config) GetQueueLength(serviceId int64) (int, error) {
	val, err := app.Models.Queue.GetQueueLength(serviceId)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}
