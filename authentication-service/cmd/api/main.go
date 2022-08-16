package main

import (
	"authentication-service/cmd/api/handler"
	"authentication-service/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const servicePort = "8081"

type Config struct {
	DB      *pgxpool.Pool
	Models  data.Models
	Handler handler.Handler
}

var counts int64

func main() {
	// Connect to db .
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
		os.Exit(1)
	}
	defer conn.Close()

	app := Config{
		DB:      conn,
		Models:  data.New(conn),
		Handler: handler.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", servicePort),
		Handler: app.routes(),
	}
	log.Printf("The authentication service runnig on port %s\n", servicePort)
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() *pgxpool.Pool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	db_url := os.Getenv("DB_URL")
	for {
		conn, err := pgxpool.Connect(context.Background(), db_url)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return conn
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
