package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("sms-jwt-secret"), nil)
}

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

	// Protected routes
	mux.Group(func(mux chi.Router) {
		mux.Use(jwtauth.Verifier(tokenAuth))
		mux.Use(jwtauth.Authenticator)

	})

	// Public routes
	mux.Group(func(mux chi.Router) {
		mux.Get("/", app.Broker)
		mux.Post("/getToken", app.getToken)
	})

	return mux
}
