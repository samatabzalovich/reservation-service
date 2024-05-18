package main

import (
	"fmt"
)



func (app *Config) GetAppointmentsForClientInInstitution(clientId int64, institutionId , employeeId string) (int, error) {
	endpoint := fmt.Sprintf("http://%s%s", app.appointmentHost, fmt.Sprintf("/completed-appointments-number/{%d}?instId=%s&employeeId=%s", clientId, institutionId, employeeId))
	var output struct{
		Count int `json:"count"`
	}
	
	err := app.sendGetRequest(endpoint, output)

	if err != nil {
		return 0, err
	}

	return output.Count, nil
}

func (app *Config) GetQueueForClientInInstitution(clientId int64, institutionId , employeeId string) (int, error) {
	endpoint := fmt.Sprintf("http://%s%s", app.queueHost, fmt.Sprintf("/queue-number/{%d}?instId=%s&employeeId=%s", clientId, institutionId, employeeId))
	var output struct{
		Count int `json:"count"`
	}
	
	err := app.sendGetRequest(endpoint, output)

	if err != nil {
		return 0, err
	}

	return output.Count, nil
}