package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"queue-managemant-service/internal/data"

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
	upgrader         websocket.Upgrader
	broadcast        chan Message
	hub              *Hub
	Redis            *redis.Client
	authServiceHost string
}

func main() {
	authHost := os.Getenv("AUTH_SERVICE")
	if authHost == "" {
		authHost = "localhost:50001"
	}
	db := connectToDB()
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	app := Config{
		port:             "8087",
		Models:           data.New(db),
		staffServiceHost: os.Getenv("staff_service_host"),
		Redis:            openRedisConn(),
		upgrader:         upgrader,
		broadcast:        make(chan Message),
		hub:              NewHub(),
		authServiceHost: authHost,
	}

	//err := app.InitNotificationSender()
	//if err != nil {
	//	log.Panic(err)
	//}
	//defer app.amqpConn.Close()
	//defer app.Ch.Close()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.routes(),
	}
	go app.hub.Run()
	log.Printf("Starting queue service on port %s\n", app.port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
