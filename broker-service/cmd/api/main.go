package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	port string
	authServiceHost string
	instHost string
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

	app := Config{
		port: "80",
		authServiceHost: authHost,
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.routes(),
	}
	log.Printf("Starting broker sevice on port %s\n", app.port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
