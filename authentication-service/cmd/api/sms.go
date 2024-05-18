package main

import (
	"context"
	"time"
)

func (authServer *AuthService) checkSmsCode(ctx context.Context, code string, number string) bool {
	val, err := authServer.Redis.Get(ctx, number).Result()
	if err != nil {
		return false
	}
	if val != code {
		return false
	}
	return true
}

func (authServer *AuthService) sendSms(ctx context.Context, number string) error {
	code := authServer.Sender.GenerateCode()
	err := authServer.Redis.Set(ctx, number, code, 5*60*time.Second).Err()
	if err != nil {
		return err
	}
	err = authServer.Sender.SendSmsCodeAPI(number, code)
	if err != nil {
		return err
	}
	return nil
}
