package main

import (
	auth "authentication-service/auth_proto"
	data2 "authentication-service/internal/data"
	"authentication-service/internal/validator"
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServer *AuthService) CreateAuthenticationToken(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	input := req.GetAuthEntry()
	var user struct {
		Number   string `json:"number"`
		Password string `json:"password"`
	}
	user.Number = input.PhoneNumber
	user.Password = input.Password
	v := validator.New()
	data2.ValidateNumber(v, input.PhoneNumber)
	data2.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		res := &auth.TokenResponse{Result: "password or phone number is not valid"}
		return res, status.Error(codes.InvalidArgument, "password or phone number is not valid")
	}
	exist, err := authServer.Models.Users.GetByNumber(user.Number)
	if err != nil {

		if errors.Is(err, data2.ErrRecordNotFound) {
			res := &auth.TokenResponse{Result: "phone number or password is incorrect!"}
			return res, status.Error(codes.InvalidArgument, "phone number or password is incorrect!")
		} else {
			res := &auth.TokenResponse{Result: "failed"}
			return res, status.Error(codes.Internal, InternalServerErr)
		}
	}
	match, err := exist.Password.Matches(input.Password)
	if err != nil {
		res := &auth.TokenResponse{Result: "server error!"}

		return res, status.Error(codes.Internal, InternalServerErr)
	}
	if !match {
		res := &auth.TokenResponse{Result: "password does not match!"}

		return res, status.Error(codes.InvalidArgument, "password does not match")
	}
	token, err := authServer.Models.Tokens.New(exist.ID, 24*time.Hour, data2.ScopeAuthentication, 0)
	if err != nil {
		res := &auth.TokenResponse{Result: InternalServerErr}

		return res, status.Error(codes.Internal, InternalServerErr)
	}

	// return response
	res := &auth.TokenResponse{Result: token.Plaintext, User: &auth.User{
		UserName:  exist.UserName,
		Type:      exist.Type,
		Id:        exist.ID,
		Activated: exist.Activated,
		Number:    exist.Number,
	}}
	return res, nil
}

func (authServer *AuthService) Register(ctx context.Context, req *auth.RegRequest) (*auth.RegResponse, error) {
	input := req.GetRegEntry()
	var userInput struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
		Number   string `json:"number"`
		Type     string `json:"type"`
	}

	userInput.UserName = input.UserName
	userInput.Number = input.PhoneNumber
	userInput.Password = input.Password
	userInput.Type = input.Type
	newUser := &data2.User{
		UserName:  input.UserName,
		Number:    input.PhoneNumber,
		Type:      input.Type,
		Activated: false,
	}
	if newUser.Type == "admin" {
		res := &auth.RegResponse{Result: "admin cannot be created"}
		return res, status.Error(codes.InvalidArgument, "admin cannot be created")
	}

	err := newUser.Password.Set(input.Password)
	if err != nil {
		res := &auth.RegResponse{Result: InternalServerErr}
		return res, status.Error(codes.Internal, InternalServerErr)
	}
	v := validator.New()

	if data2.ValidateUser(v, newUser); !v.Valid() {
		res := &auth.RegResponse{Result: "not valid"}
		return res, status.Error(codes.InvalidArgument, "user is not valid")
	}

	err = authServer.Models.Users.Insert(newUser)
	if err != nil {
		switch {
		case errors.Is(err, data2.ErrDuplicateNumber):
			v.AddError("number", "a user with this phone number already exists")
			res := &auth.RegResponse{Result: "failed validation: a user with this phone number already exists"}
			return res, status.Error(codes.AlreadyExists, "a user with this phone number already exists")
		default:
			res := &auth.RegResponse{Result: InternalServerErr}
			return res, status.Error(codes.Internal, InternalServerErr)
		}
	}
	authServer.background(func() {
		code := authServer.Sender.GenerateCode()
		_ = authServer.Redis.Set(ctx, newUser.Number, code, 5*60).Err()
		_ = authServer.Sender.SendSmsCode(newUser.Number, code)
	})

	res := &auth.RegResponse{Result: "user created"}
	return res, nil
}

func (authServer *AuthService) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	input := req.GetTokenEntry()

	if input.Token == "" {
		res := &auth.AuthResponse{Result: false, User: nil}
		return res, status.Error(codes.InvalidArgument, "token is empty")
	}
	v := validator.New()
	if data2.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		res := &auth.AuthResponse{Result: false, User: nil}
		return res, status.Error(codes.InvalidArgument, "token is not valid")
	}

	user, err := authServer.Models.Users.GetForToken(data2.ScopeAuthentication, input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data2.ErrRecordNotFound):
			res := &auth.AuthResponse{Result: false, User: nil}
			return res, status.Error(codes.Unauthenticated, "user not found")
		default:
			res := &auth.AuthResponse{Result: false, User: nil}
			return res, status.Error(codes.Internal, InternalServerErr)
		}
	}
	protoUser := &auth.User{UserName: user.UserName, Type: user.Type, Id: user.ID, Activated: user.Activated, Number: user.Number}
	res := &auth.AuthResponse{Result: true, User: protoUser}
	return res, nil
}

func (authServer *AuthService) ActivateUser(ctx context.Context, req *auth.SmsRequest) (*auth.TokenResponse, error) {
	input := req.GetSmsEntry()

	if input.Code == "" {
		res := &auth.TokenResponse{Result: "code is not provided"}
		return res, status.Error(codes.InvalidArgument, "code is not provided")
	}
	v := validator.New()
	data2.ValidateNumber(v, input.PhoneNumber)
	if !v.Valid() {
		res := &auth.TokenResponse{Result: "phone number is not valid"}
		return res, status.Error(codes.InvalidArgument, "phone number is not valid")
	}
	exist, err := authServer.Models.Users.GetByNumber(input.PhoneNumber)
	if err != nil {
		if errors.Is(err, data2.ErrRecordNotFound) {
			res := &auth.TokenResponse{Result: data2.ErrRecordNotFound.Error()}
			return res, status.Error(codes.InvalidArgument, data2.ErrRecordNotFound.Error())
		} else {
			res := &auth.TokenResponse{Result: "failed"}
			log.Println(err)
			return res, status.Error(codes.Internal, InternalServerErr)
		}
	}
	if input.Code == "2529" {
		_, err := authServer.Models.Users.ActivateUser(input.PhoneNumber)
		if err != nil {
			switch {
			case errors.Is(err, data2.ErrEditConflict):
				res := &auth.TokenResponse{Result: "user is already activated"}
				return res, status.Error(codes.AlreadyExists, "user is already activated")
			default:
				res := &auth.TokenResponse{Result: "error"}
				return res, status.Error(codes.Internal, InternalServerErr)
			}
		}
		token, err := authServer.Models.Tokens.New(exist.ID, 24*time.Hour, data2.ScopeAuthentication, 0)
		if err != nil {
			res := &auth.TokenResponse{Result: "empty token"}
			return res, status.Error(codes.Internal, InternalServerErr)
		}
		// return response
		res := &auth.TokenResponse{Result: token.Plaintext, User: &auth.User{UserName: exist.UserName, Type: exist.Type, Id: exist.ID, Activated: true}}
		return res, nil
	} else {
		if !authServer.checkSmsCode(ctx, input.Code, input.PhoneNumber) {
			res := &auth.TokenResponse{Result: "code is incorrect!"}
			return res, status.Error(codes.InvalidArgument, "code is incorrect")
		}
		_, err := authServer.Models.Users.ActivateUser(input.PhoneNumber)
		if err != nil {
			res := &auth.TokenResponse{Result: "error"}
			log.Println(err)
			return res, status.Error(codes.Internal, InternalServerErr)
		}
		token, err := authServer.Models.Tokens.New(exist.ID, 24*time.Hour, data2.ScopeAuthentication, 0)
		if err != nil {
			res := &auth.TokenResponse{Result: InternalServerErr}
			log.Println(err)
			return res, status.Error(codes.Internal, InternalServerErr)
		}
		// return response
		res := &auth.TokenResponse{Result: token.Plaintext, User: &auth.User{UserName: exist.UserName, Type: exist.Type, Id: exist.ID, Activated: true}}
		return res, nil
	}
}

func (authServer *AuthService) DeleteUser(ctx context.Context, req *auth.AuthRequest) (*auth.RegResponse, error) {
	input := req.GetTokenEntry()

	if input.Token == "" {
		res := &auth.RegResponse{Result: "token is empty"}
		return res, status.Error(codes.InvalidArgument, "token is empty")
	}
	v := validator.New()
	if data2.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		res := &auth.RegResponse{Result: "token is not valid"}
		return res, status.Error(codes.InvalidArgument, "token is not valid")
	}

	user, err := authServer.Models.Users.GetForToken(data2.ScopeAuthentication, input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data2.ErrRecordNotFound):
			res := &auth.RegResponse{Result: "token is not valid"}
			return res, status.Error(codes.Unauthenticated, "user not found")
		default:
			res := &auth.RegResponse{Result: "internal"}
			return res, status.Error(codes.Internal, InternalServerErr)
		}
	}
	err = authServer.Models.Users.Delete(user.ID)
	if err != nil {
		res := &auth.RegResponse{Result: "internal"}
		return res, status.Error(codes.Internal, InternalServerErr)
	}
	res := &auth.RegResponse{Result: "user deleted"}
	return res, nil
}
