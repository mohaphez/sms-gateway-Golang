package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sms-service/cmd/api/handler"
	"sms-service/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	servicePort = "8082"
	mongoURL    = "mongodb://127.0.0.1:27017"
)

type Config struct {
	Models  data.Models
	Handler handler.Handler
}

func main() {
	// Connect to db .
	conn, err := connectToDB()
	if err != nil {
		log.Panic(err)
	}

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = conn.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models:  data.New(conn),
		Handler: handler.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", servicePort),
		Handler: app.routes(),
	}

	log.Printf("The sms service runnig on port %s\n", servicePort)

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MOGODB_USERNAME"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error Connecting to mongo db : ", err)
		return nil, err
	}

	log.Println("Connected to mongo !")
	return conn, nil
}
