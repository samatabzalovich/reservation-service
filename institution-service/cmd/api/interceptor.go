package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (app *Config) authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s", info.FullMethod)
	protectedMethods := []string{
		"/inst.InstitutionService/CreateInstitution",
		"/inst.InstitutionService/UpdateInstitution",
		"/inst.InstitutionService/DeleteInstitution",
		"/inst.CategoryService/CreateCategory",
	}

	if isInProtectedMethods(info.FullMethod, protectedMethods) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metadata is not provided")
		}

		tokens, ok := md["authorization"]
		if !ok || len(tokens) == 0 {
			return nil, errors.New("authorization token is not provided")
		}

		token := tokens[0]
		payload, err := app.GetUserForToken(token)
		if err != nil || payload == nil {
			return nil, errors.New("invalid token")
		}
		if !payload.Activated {
			return nil, errors.New("user is not activated")
		}
		// Add the user data to the context
		ctx = *app.contextSetUser(ctx, payload)
	}

	// Continue processing the request
	return handler(ctx, req)
}

func isInProtectedMethods(method string, protectedMethods []string) bool {
	for _, protectedMethod := range protectedMethods {
		if strings.EqualFold(method, protectedMethod) {
			return true
		}
	}
	return false
}
