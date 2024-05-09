package main

import (
	"encoding/json"
	"fmt"
	"log"
	"notification-service/internal/data"
	"os"

	"github.com/streadway/amqp"
)

func (app *Config) ListenForAppointmentNotificationRequests() {
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHost == "" {
		log.Fatalf("RABBITMQ_HOST env variable not set")
	}
	conn, err := amqp.Dial(rabbitMQHost)
	if err != nil {
		log.Fatalf("Error occurred connecting to RabbitMQ. Err: %s", err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error occurred creating a channel. Err: %s", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"appointment-notification", // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)

	if err != nil {
		log.Fatalf("Error occurred creating a queue. Err: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("Error occurred consuming messages. Err: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println("appointment fetched")
			var appointment data.AppointmentData
			_ = json.Unmarshal(d.Body, &appointment)
			photoUrl, err := app.Models.Appointments.GetPhotoURL(appointment.ID)
			if err != nil {
				log.Printf("Error occurred getting photo URL. Err: %s", err)
				continue
			}
			clientDevices, err := app.Models.DeviceTokens.GetByUserID(int64(appointment.ClientID))
			if err != nil {
				log.Printf("Error occurred getting device tokens. Err: %s", err)
				continue
			}
			employeeDevices, err := app.Models.DeviceTokens.GetByEmployeeID(int64(appointment.EmployeeID))
			if err != nil {
				log.Printf("Error occurred getting device tokens. Err: %s", err)
				continue
			}
			log.Println("photo: ", photoUrl)
			for _, clientDevice := range clientDevices {
				request := RequestBody{
					Token: clientDevice.Token,
					Data: &Message{
						Title:    "Appointment approved",
						Body:     "Your appointment has been approved",
						ImageUrl: photoUrl,
					},
					Notification: &Message{
						Title:    "Appointment approved",
						Body:     "Your appointment has been approved",
						ImageUrl: photoUrl,
					},
				}
				sent := app.sendRequest(request)
				if !sent {
					log.Printf("Error occurred sending request. Err: %s", err)
					continue
				}

			}
			for _, employeeDevice := range employeeDevices {
				h, m, _ := appointment.StartTime.Time.Clock()
				
				notificationString := fmt.Sprintf("on %d.%d at %s", appointment.StartTime.Time.Day(), appointment.StartTime.Time.Month(), fmt.Sprintf("%d:%d", h, m))
				request := RequestBody{
					Token: employeeDevice.Token,
					Data: &Message{
						Title:    "Appointment approved",
						Body:     fmt.Sprintf("You have an appointment %s", notificationString),
						ImageUrl: photoUrl,
					},
					Notification: &Message{
						Title:    "Appointment approved",
						Body:     fmt.Sprintf("You have an appointment %s", notificationString),
						ImageUrl: photoUrl,
					},
				}
				sent := app.sendRequest(request)
				if !sent {
					log.Printf("Error occurred sending request. Err: %s", err)
					continue
				}
			}
		}
	}()

	<-forever
}
