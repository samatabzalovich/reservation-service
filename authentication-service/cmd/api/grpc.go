package main

import (
	auth "authentication-service/auth_proto"
	data2 "authentication-service/internal/data"
	"authentication-service/internal/sms"
	"github.com/redis/go-redis/v9"
	"sync"
)

type AuthService struct {
	auth.UnimplementedTokenServiceServer
	auth.UnimplementedRegServiceServer
	auth.UnimplementedAuthServiceServer
	auth.UnimplementedSmsServiceServer
	Models data2.Models
	wg     sync.WaitGroup
	Sender *sms.MessageService
	Redis  *redis.Client
}
