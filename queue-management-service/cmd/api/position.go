package main

import (
	"context"
	"strconv"
)

func (app *Config) GetQueueLength(serviceId int64, ctx context.Context) (int, error) {
	val, err := app.Redis.Get(ctx, strconv.FormatInt(serviceId, 10)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			app.Redis.Set(ctx, strconv.FormatInt(serviceId, 10), 0, 0)
			return 0, nil
		} else {
			return 0, err
		}

	}
	length, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return length, nil
}

func (app *Config) IncreaseQueueLength(serviceId int64, ctx context.Context) error {
	err := app.Redis.Incr(ctx, strconv.FormatInt(serviceId, 10)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) DecreaseQueueLength(serviceId int64, ctx context.Context) error {
	err := app.Redis.Decr(ctx, strconv.FormatInt(serviceId, 10)).Err()
	if err != nil {
		return err
	}
	return nil
}
