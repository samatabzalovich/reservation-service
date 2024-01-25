package main

import (
	auth "authentication-service/auth_proto"
	"authentication-service/internal/data"
	"authentication-service/internal/sms"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const grpcPort = "50001"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
	Redis  *redis.Client
}

func main() {
	log.Println("Starting authentication service")
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
		Redis:  openRedisConn(),
	}

	go func() {
		log.Fatal(http.ListenAndServe(":8084", nil))
	}()
	app.grpcListen()

}

func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	s := grpc.NewServer()
	//get url from env file
	url := os.Getenv("SMS_URL")
	//get API key from env file
	apiKey := os.Getenv("API_KEY")
	authService := &AuthService{Models: app.Models, Sender: sms.NewMessageService(url, apiKey), Redis: app.Redis}
	auth.RegisterTokenServiceServer(s, authService)
	auth.RegisterRegServiceServer(s, authService)
	auth.RegisterAuthServiceServer(s, authService)
	auth.RegisterSmsServiceServer(s, authService)
	log.Printf("gRPC Server started on port %s", grpcPort)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
