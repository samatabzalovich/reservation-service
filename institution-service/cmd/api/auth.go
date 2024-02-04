package main

import (
	"context"
	auth "institution-service/proto_files/auth_proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
type AuthPayload struct {
	UserId int64
	Type string
	Activated bool
}

func (app *Config) GetUserForToken(token string) (*AuthPayload, error) {

	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
	return &AuthPayload{
		UserId: res.User.GetId(),
		Type: res.User.GetType(),
		Activated: res.User.GetActivated(),
	}, nil
}