package main

import (
	"authentication-service/cmd/api/handler"
	"net/http"

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
	mux.Post("/token", handler.Authentication)
	mux.Post("/create-user", handler.Createuser)

	return mux
}
