package main

import (
	auth "authentication-service/auth_proto"
	employee "authentication-service/employee_proto"
	data2 "authentication-service/internal/data"
	"authentication-service/internal/sms"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	InternalServerErr = "data server error"
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

type EmployeeService struct {
	employee.UnimplementedTokenEmployeeRegisterServiceServer
	Models data2.Models
}
