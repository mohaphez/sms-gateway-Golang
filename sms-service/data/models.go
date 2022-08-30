package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client

func New(mongodb *mongo.Client) Models {
	db = mongodb
	return Models{
		SendSMS: SendSMS{},
	}
}

type Models struct {
	SendSMS SendSMS
}
