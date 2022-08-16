package handler

import (
	"authentication-service/cmd/api/utils"
	"authentication-service/data"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Handler struct{}

type Config struct {
	DB     *pgxpool.Pool
	Models data.Models
}

var app = Config{}

func New(dbpool *pgxpool.Pool) Handler {
	app = Config{
		DB:     dbpool,
		Models: data.New(dbpool),
	}
	return Handler{}
}

func Welcome(w http.ResponseWriter, r *http.Request) {

	res := utils.JsonResponse{
		Error:   false,
		Message: "Welcome to the authentication service .",
	}
	utils.WriteJSON(w, http.StatusAccepted, res)
}
