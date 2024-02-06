package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	port string
}

func main() {

	app := Config{
		port: "80",
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
