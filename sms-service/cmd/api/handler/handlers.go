package handler

import (
	"net/http"
	"sms-service/cmd/api/utils"
	"sms-service/data"

	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct{}

type Config struct {
	DB     *mongo.Client
	Models data.Models
}

var app = Config{}

func New(mongodb *mongo.Client) Handler {
	app = Config{
		DB:     mongodb,
		Models: data.New(mongodb),
	}
	return Handler{}
}

func Welcome(w http.ResponseWriter, r *http.Request) {

	res := utils.JsonResponse{
		Error:   false,
		Message: "Welcome to the sms service .",
	}
	utils.WriteJSON(w, http.StatusAccepted, res)
}
