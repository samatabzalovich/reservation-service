package main

import (
	employee "authentication-service/employee_proto"
	"authentication-service/internal/data"
	"context"
	"time"
)
func (employeeService *EmployeeService) RegisterEmployee(ctx context.Context, req *employee.TokenEmployeeRegisterRequest) (*employee.TokenEmployeeRegisterResponse, error) {
	instId := req.GetInstitutionId()
	ownerId := req.GetOwnerId()
	employeeRegToken, err := employeeService.Models.Tokens.New(ownerId, 5*time.Minute, data.ScopeEmployeeReg, instId)
	if err != nil {
		res := &employee.TokenEmployeeRegisterResponse{Result: "server error"}
		return res, err
	}

	// return response
	res := &employee.TokenEmployeeRegisterResponse{Token: employeeRegToken.Plaintext, Result: "employee registration token created"}
	return res, nil
}