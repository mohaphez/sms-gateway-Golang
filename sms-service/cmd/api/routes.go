package main

import (
	"net/http"
	"sms-service/cmd/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {

	mux := chi.NewRouter()

	// Set cors options for who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// routes middleware
	mux.Use(middleware.Heartbeat("/ping"))

	// routes list ......
	mux.Get("/", handler.Welcome)
	mux.Post("/send-sms", handler.SendSMS)
	mux.Post("/send-sms-p2p", handler.SendSMSArray)
	mux.Post("/get-sms-status", handler.SmsStatus)
	mux.Post("/get-sms-list", handler.SmsList)
	return mux
}
