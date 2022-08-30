package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SendSMS struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Sender         string             `bson:"sender" json:"sender"`
	SenderNumber   string             `bson:"sender_number" json:"sender_number"`
	Receptor       string             `bson:"receptor" json:"receptor"`
	SendType       int                `bson:"send_type" json:"send_type"`
	SmsCount       int                `bson:"sms_count" json:"sms_count"`
	Message        string             `bson:"message" json:"message"`
	Lang           string             `bson:"lang" json:"lang"`
	Status         int                `bson:"status" json:"status"`
	StatusText     string             `bson:"status_text" json:"status_text"`
	Identity       string             `bson:"identity,omitempty" json:"identity,omitempty"`
	DeliveryStatus string             `bson:"delivery_status,omitempty" json:"delivery_status,omitempty"`
	BatchId        string             `bson:"batchid,omitempty" json:"batchid,omitempty"`
	LocalId        string             `bson:"localid,omitempty" json:"localid,omitempty"`
	User           string             `bson:"user,omitempty" json:"user,omitempty"`
	Error          string             `bson:"error,omitempty" json:"error,omitempty"`
	Date           int64              `bson:"date,omitempty" json:"date,omitempty"`
	SendTime       time.Time          `bson:"send_time,omitempty" json:"send_time,omitempty"`
	ReceiveTime    time.Time          `bson:"receive_time,omitempty" json:"receive_time,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

func (s *SendSMS) Insert(sms SendSMS) (*mongo.InsertOneResult, SendSMS, error) {
	collection := db.Database("sms-service").Collection("send")

	messageid, err := collection.InsertOne(context.TODO(), SendSMS{
		Sender:         sms.Sender,
		SenderNumber:   sms.SenderNumber,
		Receptor:       sms.Receptor,
		SendType:       sms.SendType,
		SmsCount:       sms.SmsCount,
		Message:        sms.Message,
		Status:         sms.Status,
		StatusText:     sms.StatusText,
		Identity:       sms.Identity,
		DeliveryStatus: sms.DeliveryStatus,
		BatchId:        sms.BatchId,
		Lang:           sms.Lang,
		User:           sms.User,
		Error:          sms.Error,
		Date:           sms.Date,
		SendTime:       sms.SendTime,
		ReceiveTime:    sms.ReceiveTime,
		CreatedAt:      sms.CreatedAt,
		UpdatedAt:      sms.UpdatedAt,
	})

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return nil, sms, err
	}

	return messageid, sms, nil
}

func (s *SendSMS) Update(id interface{}, sms SendSMS) error {
	collection := db.Database("sms-service").Collection("send")
	update := bson.M{
		"$set": bson.M{
			"sender":          sms.Sender,
			"sender_number":   sms.SenderNumber,
			"receptor":        sms.Receptor,
			"send_type":       sms.SendType,
			"sms_count":       sms.SmsCount,
			"message":         sms.Message,
			"status":          sms.Status,
			"status_text":     sms.StatusText,
			"identity":        sms.Identity,
			"delivery_status": sms.DeliveryStatus,
			"batchId":         sms.BatchId,
			"lang":            sms.Lang,
			"user":            sms.User,
			"send_time":       sms.SendTime,
			"receive_time":    sms.ReceiveTime,
			"error":           sms.Error,
			"createdAt":       sms.CreatedAt,
			"updatedAt":       sms.UpdatedAt,
		},
	}
	_, err := collection.UpdateByID(context.TODO(), id.(primitive.ObjectID), update)

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return nil
	}

	return nil
}
