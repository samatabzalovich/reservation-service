package main

import (
	"fmt"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
)

const grpcPort = "50002" 

var counts int64

type Config struct {
	Models          data.Models
	authServiceHost string
}

func main() {
	log.Println("Starting institution service")
	authHost := os.Getenv("AUTH_SERVICE")
	if authHost == "" {
		log.Fatal("AUTH_SERVICE env variable is not set")
	}
	dbConn := connectToDB()
	if dbConn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	app := &Config{
		Models:          data.New(dbConn),
		authServiceHost: authHost,
	}
	http.HandleFunc("/health", app.HealthCheck)

	// start http server
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8091"
		}
		log.Println("Starting HTTP health check server on port ", port)
		
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}()
	app.grpcListen()
}

func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(app.authUnaryInterceptor),
	)
	institutionService := &InstitutionService{Models: app.Models}
	categoryService := &CategoryService{Models: app.Models}
	inst.RegisterInstitutionServiceServer(s, institutionService)
	inst.RegisterCategoryServiceServer(s, categoryService)
	log.Printf("gRPC Server started on port %s", grpcPort)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
