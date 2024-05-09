package main

import (
	"appointment-service/internal/data"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func (app *Config) InitNotificationSender() error {
	rabbitMQHOST := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHOST == "" {
		log.Fatal("RABBITMQ_HOST env variable is not set")
	}
	conn, err := amqp.Dial(rabbitMQHOST)
	if err != nil {
		return err
	}
	
	log.Println("Connected to RabbitMQ")
	// Open a channel for communication
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	
	// Declare a durable queue
	q, err := ch.QueueDeclare(
		"appointment-notification", // Queue name
		true,       // Durable (persists even if RabbitMQ restarts)
		false,      // Not exclusive
		false,      // Not auto-deleted
		false,      // No wait
		nil,        // Additional arguments
	)
	if err != nil {
		return err
	}
	app.amqpConn = conn
	app.Ch = ch
	app.Queue = q
	return nil
}

func (app *Config) SendAppointmentNotification(appointment data.Appointment) error {
	body, err := json.Marshal(appointment)
	if err != nil {
		return err
	}
	return app.Ch.Publish(
		"",           // Exchange
		app.Queue.Name, // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
