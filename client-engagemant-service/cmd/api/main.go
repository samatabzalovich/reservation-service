package main

import (
	"client-engagemant-service/internal/data"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/streadway/amqp"
)

type Config struct {
	port             string
	Models           data.Models
	staffServiceHost string
	Ch               *amqp.Channel
	Queue            amqp.Queue
	amqpConn         *amqp.Connection
	appointmentHost string
	appointmentBasePath string
	queueBasePath string
	queueHost string
}

func main() {
	db := connectToDB()
	appoinstmentHost := os.Getenv("APPOINTMENT_HOST")
	if appoinstmentHost == "" {
		appoinstmentHost = "localhost:8083"
	}
	queueHost := os.Getenv("QUEUE_HOST")
	if queueHost == "" {
		queueHost = "localhost:8087"
	}
	staffServie := os.Getenv("STAFF-HOST")
	if staffServie == "" {
		staffServie = "localhost:8082"
	}
	appointmentBasePath := os.Getenv("APPOINTMENT_BASE_PATH")
	if appointmentBasePath == "" {
		appointmentBasePath = "/appointment-service"
	}
	queueBasePath := os.Getenv("QUEUE_BASE_PATH")
	if queueBasePath == "" {
		queueBasePath = "/queue-service"
	}
	app := Config{
		port:             "8088",
		Models:           data.New(db),
		staffServiceHost:staffServie,
		appointmentHost: appoinstmentHost,
		queueHost: queueHost,
		appointmentBasePath: appoinstmentHost,
		queueBasePath: queueBasePath,
	}

	err := app.InitNotificationSender()
	if err != nil {
		log.Panic(err)
	}
	defer app.amqpConn.Close()
	defer app.Ch.Close()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.routes(),
	}
	log.Printf("Starting employee sevice on port %s\n", app.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
