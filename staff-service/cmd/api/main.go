package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"staff-service/internal/data"

	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
)

type Config struct {
	port            string
	upgrader        websocket.Upgrader
	broadcast       chan Message
	hub             *Hub
	Models          data.Models
	authServiceHost string
	instServiceHost string
}

func main() {
	authHost := os.Getenv("AUTH_SERVICE")
	if authHost == "" {
		log.Fatal("AUTH_SERVICE env variable is not set")
	}
	instHost := os.Getenv("INSTITUTION_SERVICE")
	if instHost == "" {
		log.Fatal("INSTITUTION_SERVICE env variable is not set")
	}
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db := connectToDB()
	app := Config{
		port:            "80",
		upgrader:        upgrader,
		broadcast:       make(chan Message),
		hub:             NewHub(),
		Models:          data.New(db),
		authServiceHost: authHost,
		instServiceHost: instHost,
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.routes(),
	}
	go app.hub.Run()
	log.Printf("Starting employee sevice on port %s\n", app.port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
