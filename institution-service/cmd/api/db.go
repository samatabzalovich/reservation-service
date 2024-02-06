package main

import (
	"database/sql"
	// "github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

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

// func openRedisConn() *redis.Client {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "redis:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	return client
// }
