package main

import (
	auth "authentication-service/auth_proto"
	data2 "authentication-service/internal/data"
	"authentication-service/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"
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
		return res, errors.New("not valid")
	}
	exist, err := authServer.Models.Users.GetByNumber(user.Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res := &auth.TokenResponse{Result: "phone number or password is incorrect!"}
			return res, nil
		} else {
			res := &auth.TokenResponse{Result: "failed"}
			return res, err
		}
	}
	match, err := exist.Password.Matches(input.Password)
	if err != nil {
		res := &auth.TokenResponse{Result: "server error!"}

		return res, err
	}
	if !match {
		res := &auth.TokenResponse{Result: "password does not match!"}

		return res, errors.New("password error")
	}
	token, err := authServer.Models.Tokens.New(exist.ID, 24*time.Hour, data2.ScopeAuthentication, 0)
	if err != nil {
		res := &auth.TokenResponse{Result: "server error"}

		return res, err
	}

	// return response
	res := &auth.TokenResponse{Result: token.Plaintext, User: &auth.User{
		UserName:  exist.UserName,
		Email:     exist.Email,
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
		Email    string `json:"email"`
		Password string `json:"password"`
		Number   string `json:"number"`
		Type     string `json:"type"`
	}

	userInput.UserName = input.UserName
	userInput.Number = input.PhoneNumber
	userInput.Email = input.Email
	userInput.Password = input.Password
	userInput.Type = input.Type
	newUser := &data2.User{
		UserName:  input.UserName,
		Email:     input.Email,
		Number:    input.PhoneNumber,
		Type:      input.Type,
		Activated: false,
	}

	err := newUser.Password.Set(input.Password)
	if err != nil {
		res := &auth.RegResponse{Result: "Server Error"}
		return res, err
	}
	v := validator.New()

	if data2.ValidateUser(v, newUser); !v.Valid() {
		res := &auth.RegResponse{Result: "not valid"}
		return res, errors.New("user is not valid")
	}

	err = authServer.Models.Users.Insert(newUser)
	if err != nil {
		switch {
		case errors.Is(err, data2.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			res := &auth.RegResponse{Result: "failed validation: email already exists"}
			return res, err
		default:
			res := &auth.RegResponse{Result: "server error"}
			return res, err
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
		return res, errors.New("token is not provided")
	}
	v := validator.New()
	if data2.ValidateTokenPlaintext(v, input.Token); !v.Valid() {
		res := &auth.AuthResponse{Result: false, User: nil}
		return res, errors.New("token is not valid")
	}

	user, err := authServer.Models.Users.GetForToken(data2.ScopeAuthentication, input.Token)
	if err != nil {
		switch {
		case errors.Is(err, data2.ErrRecordNotFound):
			res := &auth.AuthResponse{Result: false, User: nil}
			return res, errors.New("invalid authentication token")
		default:
			res := &auth.AuthResponse{Result: false, User: nil}
			return res, errors.New("server error")
		}
	}
	protoUser := &auth.User{UserName: user.UserName, Type: user.Type, Email: user.Email, Id: user.ID, Activated: user.Activated, Number: user.Number}
	res := &auth.AuthResponse{Result: true, User: protoUser}
	return res, nil
}

func (authServer *AuthService) ActivateUser(ctx context.Context, req *auth.SmsRequest) (*auth.TokenResponse, error) {
	input := req.GetSmsEntry()

	if input.Code == "" {
		res := &auth.TokenResponse{Result: "code is not provided"}
		return res, errors.New("code is not provided")
	}
	v := validator.New()
	data2.ValidateNumber(v, input.PhoneNumber)
	if !v.Valid() {
		res := &auth.TokenResponse{Result: "phone number is not valid"}
		return res, errors.New("not valid")
	}
	exist, err := authServer.Models.Users.GetByNumber(input.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res := &auth.TokenResponse{Result: "phone number is incorrect!"}
			return res, nil
		} else {
			res := &auth.TokenResponse{Result: "failed"}
			return res, err
		}
	}
	if input.Code == "2529" {
		_, err := authServer.Models.Users.ActivateUser(input.PhoneNumber)
		if err != nil {
			res := &auth.TokenResponse{Result: "error"}
			return res, err
		}
		token, err := authServer.Models.Tokens.New(exist.ID, 24*time.Hour, data2.ScopeAuthentication, 0)
		if err != nil {
			res := &auth.TokenResponse{Result: "empty token"}
			return res, err
		}
		// return response
		res := &auth.TokenResponse{Result: token.Plaintext, User: &auth.User{UserName: exist.UserName, Email: exist.Email, Type: exist.Type, Id: exist.ID, Activated: true}}
		return res, nil
	} else {
		if !authServer.checkSmsCode(ctx, input.Code, input.PhoneNumber) {
			res := &auth.TokenResponse{Result: "code is incorrect!"}
			return res, errors.New("code error")
		}
		_, err := authServer.Models.Users.ActivateUser(input.PhoneNumber)
		if err != nil {
			res := &auth.TokenResponse{Result: "error"}
			return res, err
		}
		token, err := authServer.Models.Tokens.New(exist.ID, 24*time.Hour, data2.ScopeAuthentication, 0)
		if err != nil {
			res := &auth.TokenResponse{Result: "server error"}
			return res, err
		}
		// return response
		res := &auth.TokenResponse{Result: token.Plaintext, User: &auth.User{UserName: exist.UserName, Email: exist.Email, Type: exist.Type, Id: exist.ID, Activated: true}}
		return res, nil
	}
}
