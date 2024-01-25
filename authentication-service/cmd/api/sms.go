package main

import "context"

func (authServer *AuthService) checkSmsCode(ctx context.Context, code string, number string) bool {
	val, err := authServer.Redis.Get(ctx, code).Result()
	if err != nil {
		return false
	}
	if val != number {
		return false
	}
	return true
}
