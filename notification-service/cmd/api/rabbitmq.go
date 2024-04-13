package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func (app *Config) ListenForNotificationRequests() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
        "notification", // name
        true,         // durable
        false,        // delete when unused
        false,        // exclusive
        false,        // no-wait
        nil,          // arguments
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
			var payload RequestBody
			_ = json.Unmarshal(d.Body, &payload)

			app.sendRequest(payload)
		}
	}()

	<-forever
}