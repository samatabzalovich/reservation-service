package main

import (
	"context"
	auth "staff-service/proto_files/auth_proto"
	employee "staff-service/proto_files/employee_proto"
	inst "staff-service/proto_files/institution_proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) AuthenticateViaGrpc(token string) (*User, error) {
	conn, err := grpc.Dial("authentication-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {

		return nil, err
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Authenticate(ctx, &auth.AuthRequest{
		TokenEntry: &auth.Token{
			Token: token,
		},
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        res.User.Id,
		Activated: res.User.Activated,
	}, nil
}

func (app *Config) SetEmployeeRegTokenViaGrpc(ownerId, instId int64) (string, error) {
	conn, err := grpc.Dial("authentication-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := employee.NewTokenEmployeeRegisterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.RegisterEmployee(ctx, &employee.TokenEmployeeRegisterRequest{
		OwnerId:       ownerId,
		InstitutionId: instId,
	})
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (app *Config) GetInstitutionForToken(token string) (*inst.Institution, error) {
	conn, err := grpc.Dial("institution-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetForToken(ctx, &inst.GetInstForTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
