package main

import (
	"fmt"
	"log"
	data "notification-service/internal/data"
)

func (app *Config) sendAppointmentNotification() {
	// get appointments from the database
	appointments, err := app.getUpcomingAppointments()
	if err != nil {
		log.Fatalf("Error occurred getting upcoming appointments. Err: %s", err)
	}

	// for each appointment, send a notification
	for _, appointment := range appointments {
		// create a notification request
		requestBody := RequestBody{
			Token: appointment.DeviceToken,
			Notification: &Message{
				Title:    "Appointment Reminder",
				Body:     fmt.Sprintf("Hi %s! You have an appointment with %s", appointment.ClientName, appointment.EmployeeName),
				ImageUrl: appointment.PhotoUrl,
			},
			Data: &Message{
				Title:    "Appointment Reminder",
				Body:     "You have an appointment with " + appointment.EmployeeName,
				ImageUrl: appointment.PhotoUrl,
			},
		}

		// send the request
		isSent := app.sendRequest(requestBody)

		if isSent {
			app.Models.Appointments.MarkAsNotified(appointment.ID)
		} else {
			app.Models.DeviceTokens.DeleteByToken(appointment.DeviceToken)
		}
	}
}

func (app *Config) getUpcomingAppointments() ([]*data.Appointment, error) {
	appointments, err := app.Models.Appointments.GetUpcomingAppointments()

	if err != nil {
		return nil, err
	}

	return appointments, nil
}
