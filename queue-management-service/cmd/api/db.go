package main

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"os"

	// "github.com/redis/go-redis/v9"
	"log"
	"time"
)

var counts int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=190704Samat dbname=reserveHub sslmode=disable timezone=UTC connect_timeout=5"
	}
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openRedisConn() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:      host+ ":" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
