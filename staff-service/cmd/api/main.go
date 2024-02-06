package main

import (
	"fmt"
	"log"
	"net/http"
	"staff-service/internal/data"

	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	port string
	upgrader websocket.Upgrader
	broadcast chan Message
	hub *Hub
	Models data.Models
}

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}	
	db := connectToDB()
	app := Config{
		port: "80",
		upgrader: upgrader,
		broadcast: make(chan Message),
		hub: NewHub(),
		Models: data.New(db),
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