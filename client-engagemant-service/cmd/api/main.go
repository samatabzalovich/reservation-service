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
	port                string
	Models              data.Models
	staffServiceHost    string
	Ch                  *amqp.Channel
	Queue               amqp.Queue
	amqpConn            *amqp.Connection
	appointmentHost     string
	appointmentBasePath string
	queueBasePath       string
	queueHost           string
	authHost            string
}

func main() {
	db := connectToDB()
	appoinstmentHost := os.Getenv("APPOINTMENT_HOST")
	if appoinstmentHost == "" {
		appoinstmentHost = "Reserve-hub-lb-1596239107.us-east-1.elb.amazonaws.com"
	}
	queueHost := os.Getenv("QUEUE_HOST")
	if queueHost == "" {
		queueHost = "Reserve-hub-lb-1596239107.us-east-1.elb.amazonaws.com"
	}
	staffServie := os.Getenv("STAFF-HOST")
	if staffServie == "" {
		staffServie = "Reserve-hub-lb-1596239107.us-east-1.elb.amazonaws.com/staff-service/"
	}
	appointmentBasePath := os.Getenv("APPOINTMENT_BASE_PATH")
	if appointmentBasePath == "" {
		appointmentBasePath = "/appointment-service-endpoint"
	}
	queueBasePath := os.Getenv("QUEUE_BASE_PATH")
	if queueBasePath == "" {
		queueBasePath = "/queue-service"
	}
	authHost := os.Getenv("AUTH_HOST")
	if authHost == "" {
		authHost = "localhost"
	}
	app := Config{
		port:                "8088",
		Models:              data.New(db),
		staffServiceHost:    staffServie,
		appointmentHost:     appoinstmentHost,
		queueHost:           queueHost,
		appointmentBasePath: appoinstmentHost,
		queueBasePath:       queueBasePath,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.routes(),
	}
	log.Printf("Starting analytics sevice on port %s\n", app.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
