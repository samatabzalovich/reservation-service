package main

import (
	auth "appointment-service/proto_files/auth_proto"
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *Config) AuthenticateViaGrpc(token string) (*User, error) {
	authDest := os.Getenv("AUTH_SERVICE")
	if authDest == "" {
		log.Fatal("AUTH_SERVICE env variable is not set")
	}
	conn, err := grpc.Dial(authDest, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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