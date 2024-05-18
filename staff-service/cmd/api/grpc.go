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

func (app *Config) SetEmployeeRegTokenViaGrpc(ownerId, instId int64) (string, error) {
	conn, err := grpc.Dial(app.authServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
	conn, err := grpc.Dial(app.instServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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

func (app *Config) GetInstitutions(token string) (*inst.Institution, error) {
	conn, err := grpc.Dial(app.instServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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

func (app *Config) GetInstitution(instId int64) (*inst.Institution, error) {
	conn, err := grpc.Dial(app.instServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetInstitution(ctx, &inst.GetInstitutionsByIdRequest{
		Id: instId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (app *Config) GetInstitutionForEmployee(employeeId int64) (*inst.Institution, error) {
	conn, err := grpc.Dial(app.instServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetInstitutionForEmployee(ctx, &inst.GetInstitutionsByIdRequest{
		Id: employeeId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (app *Config) GetInstitutionsForUserEmployee(userId int64) (*inst.InstitutionsResponse, error) {
	conn, err := grpc.Dial(app.instServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := inst.NewInstitutionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetInstitutionsForUserEmployee(ctx, &inst.GetInstitutionsByIdRequest{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}