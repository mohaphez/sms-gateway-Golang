package main

import (
	"fmt"
	"log"
	"net/http"
)

const servicePort = "8080"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Start broker service on port %s \n", servicePort)

	// Define http server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", servicePort),
		Handler: app.routes(),
	}

	// Start the borker server
	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
