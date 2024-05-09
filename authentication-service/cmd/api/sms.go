package main

import (
	"context"
	"log"
)

func (authServer *AuthService) checkSmsCode(ctx context.Context, code string, number string) bool {
	val, err := authServer.Redis.Get(ctx, code).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	if val != number {
		return false
	}
	return true
}
