package main

import (
	"fmt"
	data "institution-service/internal/data"
	inst "institution-service/proto_files/institution_proto"
	"log"
	"net"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const grpcPort = "50002" // TODO: change to 50001 when production

var counts int64

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting institution service")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConn := connectToDB()
	if dbConn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	app := &Config{
		Models: data.New(dbConn),
	}
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
