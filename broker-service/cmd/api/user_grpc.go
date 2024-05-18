package main

import (
	"broker-service/proto_files/auth"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

func (app *Config) CreateTokenViaGRpc(w http.ResponseWriter, r RequestPayload) {

	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewTokenServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.CreateAuthenticationToken(ctx, &auth.TokenRequest{
		AuthEntry: &auth.Auth{
			PhoneNumber: r.Auth.PhoneNumber,
			Password:    r.Auth.Password,
		},
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = res.Result

	app.writeJSON(w, http.StatusAccepted,
		map[string]any{
			"token": payload.Message,
			"error": payload.Error,
			"user":  res.User,
		})
}

func (app *Config) RegViaGRpc(w http.ResponseWriter, r RequestPayload) {
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewRegServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Register(ctx, &auth.RegRequest{
		RegEntry: &auth.Reg{
			Email:       r.Reg.Email,
			Password:    r.Reg.Password,
			UserName:    r.Reg.UserName,
			PhoneNumber: r.Reg.PhoneNumber,
			Type:        r.Reg.Type,
			Activated:   false,
		},
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = res.Result

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) AuthViaGRpc(w http.ResponseWriter, r RequestPayload) {
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Authenticate(ctx, &auth.AuthRequest{
		TokenEntry: &auth.Token{
			Token: r.Token.Bearer,
		},
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, map[string]any{"exist": res.Result, "user": res.User, "error": false})
}

func (app *Config) VerifySmsViaGRpc(w http.ResponseWriter, r RequestPayload) {
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()
	c := auth.NewSmsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := c.ActivateUser(ctx, &auth.SmsRequest{
		SmsEntry: &auth.Sms{
			PhoneNumber: r.Sms.PhoneNumber,
			Code:        r.Sms.Code,
		},
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, map[string]any{"token": res.Result, "user": res.User, "error": false})
}



func (app *Config) DeleteAccountViaGRpc(w http.ResponseWriter, r RequestPayload) {
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.DeleteUser(ctx, &auth.AuthRequest{
		TokenEntry: &auth.Token{
			Token: r.Token.Bearer,
		},
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, map[string]any{"message": res.Result, "error": false})
}


func (app *Config) SendSmsViaGRpc(w http.ResponseWriter, r RequestPayload) {
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewSmsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.SendCode(ctx, &auth.SmsRequest{
		SmsEntry: &auth.Sms{
			PhoneNumber: r.Sms.PhoneNumber,
		},
	})
	if err != nil {
		app.rpcErrorJson(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, map[string]any{"message": res.Result, "error": false})
}