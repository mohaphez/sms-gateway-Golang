package main

import (
	"fmt"
	"log"
	"net/http"
)

const brokerPort = "8080"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Start broker service on port %s \n", brokerPort)

	// Define http server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", brokerPort),
		Handler: app.routes(),
	}

	// Start the borker server
	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
