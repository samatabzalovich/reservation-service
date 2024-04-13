package main

import (
	employee "authentication-service/employee_proto"
	"authentication-service/internal/data"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (employeeService *EmployeeService) RegisterEmployee(ctx context.Context, req *employee.TokenEmployeeRegisterRequest) (*employee.TokenEmployeeRegisterResponse, error) {
	instId := req.GetInstitutionId()
	ownerId := req.GetOwnerId()
	employeeRegToken, err := employeeService.Models.Tokens.New(ownerId, 5*time.Minute, data.ScopeEmployeeReg, instId)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			return nil, status.Error(codes.NotFound, data.ErrRecordNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, "data server error")
		}
	}

	// return response
	res := &employee.TokenEmployeeRegisterResponse{Token: employeeRegToken.Plaintext, Result: "employee registration token created"}
	return res, nil
}
