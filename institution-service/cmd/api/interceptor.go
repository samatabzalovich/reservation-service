package main

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (app *Config) authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	protectedMethods := []string{
		"/inst.InstitutionService/CreateInstitution",
		"/inst.InstitutionService/UpdateInstitution",
		"/inst.InstitutionService/DeleteInstitution",
		"/inst.CategoryService/CreateCategory",
		"/inst.CategoryService/UpdateCategory",
		"/inst.CategoryService/DeleteCategory",
	}

	if isInProtectedMethods(info.FullMethod, protectedMethods) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "metadata is not provided")
		}

		tokens, ok := md["authorization"]
		if !ok || len(tokens) == 0 {
			return nil, status.Error(codes.PermissionDenied, "token is not provided")
		}

		token := tokens[0]
		payload, err := app.GetUserForToken(token)
		if err != nil || payload == nil {
			return nil, status.Error(codes.PermissionDenied, "invalid token")
		}
		if !payload.Activated {
			return nil, status.Error(codes.PermissionDenied, "user is not activated")
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
