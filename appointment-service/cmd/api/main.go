package main

import (
	"appointment-service/internal/data"
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
	port string
	Models data.Models
	staffServiceHost string
	Ch *amqp.Channel
	Queue amqp.Queue
	amqpConn *amqp.Connection
}

func main() {
	db := connectToDB()
	staffServiceHost := os.Getenv("STAFF_SERVICE_HOST_WITHSLASHATTHEEND")
	if staffServiceHost == "" {
		log.Fatal("STAFF_SERVICE_HOST_WITHSLASHATTHEEND env variable is not set")
	}
	app := Config{
		port: "8083",
		Models: data.New(db),
		staffServiceHost: staffServiceHost,
	}
	
	err  := app.InitNotificationSender()
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