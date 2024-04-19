package main

import (
	"fmt"
	"log"
	"net/http"
	data "notification-service/internal/data"
	"os"

	"github.com/robfig/cron/v3"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
)

type RequestBody struct {
	Token        string   `json:"token"`
	Notification *Message `json:"notification"`
	Data         *Message `json:"data"`
}

type Message struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageUrl string `json:"imageUrl"`
}

type Config struct {
	NotificationServiceHost string
	Models                  data.Models
	port string
}

func main() {
	db := connectToDB()
	//get host from env
	app := &Config{
		NotificationServiceHost: os.Getenv("host"),
		Models:                  data.New(db),
		port: "80",
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.routes(),
	}
	
	app.startCronJob()
	go app.ListenForAppointmentNotificationRequests()
	log.Printf("Starting notification sevice on port %s\n", app.port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) startCronJob() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", app.sendAppointmentNotification)
	if err != nil {
		log.Fatalf("Error adding cron job. Err: %s", err)
	}
	c.Start()
}
