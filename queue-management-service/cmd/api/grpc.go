package main

import (
	"context"
	auth "queue-managemant-service/proto_files/auth_proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) AuthenticateViaGrpc(token string) (*User, error) {
	if (app.authServiceHost == "") {
		app.authServiceHost = "localhost:50001"
	}
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
		Type:      res.User.Type,
	}, nil
}
