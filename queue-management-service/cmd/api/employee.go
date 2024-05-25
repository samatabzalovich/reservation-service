package main

func (app *Config) GetEmployeeForServiceAndUserID(serviceId, employeeUserId int64) (int64, error) {
	employee, err := app.Models.Employee.GetEmployeeForServiceAndUserID(serviceId, employeeUserId)
	if err != nil {

		return 0, err
	}

	return employee, nil
}
