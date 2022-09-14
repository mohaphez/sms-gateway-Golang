package data

import (
	"context"
	"log"
	"sms-service/cmd/api/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type SMSStatus struct {
	Sender         string    `bson:"sender" json:"sender"`
	SenderNumber   string    `bson:"sender_number" json:"sender_number"`
	Receptor       string    `bson:"receptor" json:"receptor"`
	SmsCount       int       `bson:"sms_count" json:"sms_count"`
	Message        string    `bson:"message" json:"message"`
	Lang           string    `bson:"lang" json:"lang"`
	Status         int       `bson:"status" json:"status"`
	StatusText     string    `bson:"status_text" json:"status_text"`
	DeliveryStatus string    `bson:"delivery_status,omitempty" json:"delivery_status,omitempty"`
	BatchId        string    `bson:"batchid,omitempty" json:"batchid,omitempty"`
	LocalId        string    `bson:"localid,omitempty" json:"localid,omitempty"`
	Error          string    `bson:"error,omitempty" json:"error,omitempty"`
	Date           int64     `bson:"date,omitempty" json:"date,omitempty"`
	SendTime       time.Time `bson:"send_time,omitempty" json:"send_time,omitempty"`
	ReceiveTime    time.Time `bson:"receive_time,omitempty" json:"receive_time,omitempty"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

func (s *SendSMS) Insert(sms SendSMS) (*mongo.InsertOneResult, SendSMS, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := db.Database("sms-service").Collection("send")

	messageid, err := collection.InsertOne(ctx, SendSMS{
		Sender:         sms.Sender,
		SenderNumber:   sms.SenderNumber,
		Receptor:       utils.TrimLastChars(sms.Receptor, 10),
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := db.Database("sms-service").Collection("send")
	update := bson.M{
		"$set": bson.M{
			"sender":          sms.Sender,
			"sender_number":   sms.SenderNumber,
			"receptor":        utils.TrimLastChars(sms.Receptor, 10),
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
	_, err := collection.UpdateByID(ctx, id.(primitive.ObjectID), update)

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return nil
	}

	return nil
}

func GetByBatchId(messageIds []string) (*[]SMSStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := db.Database("sms-service").Collection("send")

	res, err := collection.Find(ctx, bson.M{"batchid": bson.M{"$in": messageIds}})
	if err != nil {
		return nil, err
	}
	var messages []SMSStatus
	err = res.All(ctx, &messages)
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

func GetByFilter(message string, receptor []string, sender []string, senderNumber []string, limit int64, offset int64, sort string) (*[]SMSStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := db.Database("sms-service").Collection("send")
	opts := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
	}
	sortint := 1
	if sort == "desc" {
		sortint = -1
	}
	opts.SetSort(bson.M{"_id": sortint})

	query := bson.M{}

	if len(message) > 0 {
		query["message"] = bson.M{"$regex": primitive.Regex{Pattern: message + ".*", Options: "i"}}
	}

	if receptor != nil {
		if len(receptor) > 0 {
			query["receptor"] = bson.M{"$in": receptor}
		}
	}

	if senderNumber != nil {
		if len(senderNumber) > 0 {
			query["sender_number"] = bson.M{"$in": senderNumber}
		}
	}

	if sender != nil {
		if len(sender) > 0 {
			query["sender"] = bson.M{"$in": sender}
		}
	}

	res, err := collection.Find(ctx, query, &opts)
	if err != nil {
		return nil, err
	}
	var messages []SMSStatus
	err = res.All(ctx, &messages)
	if err != nil {
		return nil, err
	}
	return &messages, nil
}
