package handler

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"sms-service/cmd/api/provider"
	"sms-service/cmd/api/utils"
	"sms-service/data"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SendSmsResponse struct {
	Sender       string `bson:"sender" json:"sender"`
	SenderNumber string `bson:"sender_number" json:"sender_number"`
	Receptor     string `bson:"receptor" json:"receptor"`
	SmsCount     int    `bson:"sms_count" json:"sms_count"`
	Lang         string `bson:"lang" json:"lang"`
	Message      string `bson:"message" json:"message"`
	Status       int    `bson:"status" json:"status"`
	StatusText   string `bson:"status_text" json:"status_text"`
	BatchId      string `bson:"batchid" json:"batchid"`
	LocalId      string `bson:"localid,omitempty" json:"localid,omitempty"`
	Date         int64  `bson:"date" json:"date"`
}

func SendSMS(w http.ResponseWriter, r *http.Request) {
	mobileRegex, _ := regexp.Compile(`^(?:98|\+98|0098|0)?9[0-9]{9}$`)
	var SMSArray []data.SendSMS
	var responsePayload []SendSmsResponse
	var requestPayload struct {
		Message      string   `json:"message"`
		Receptor     []string `json:"receptor"`
		Sender       string   `json:"sender"`
		SenderNumber string   `json:"sender_number"`
		LocalId      string   `json:"localid,omitempty"`
	}

	// Check the request payload
	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check that the request payload are not null
	if len(requestPayload.Receptor) == 0 || len(strings.TrimSpace(requestPayload.Message)) == 0 || len(strings.TrimSpace(requestPayload.Sender)) == 0 || len(strings.TrimSpace(requestPayload.SenderNumber)) == 0 {
		utils.ErrorJSON(w, errors.New("invalid parameter ! parameters could not be null"), http.StatusBadRequest)
		return
	}

	//  Check the receptor array length are not more than 100
	if len(requestPayload.Receptor) > 100 {
		utils.ErrorJSON(w, errors.New("receptors array not be more than 100 number"), http.StatusBadRequest)
		return
	}

	smsCount, smsLanguage, err := utils.SmsCount(requestPayload.Message)
	if err != nil {
		utils.ErrorJSON(w, errors.New("message length is too high"), http.StatusBadRequest)
		return
	}

	for _, receptor := range requestPayload.Receptor {
		if mobileRegex.MatchString(receptor) {
			SMSRow := data.SendSMS{
				Sender:       requestPayload.Sender,
				SenderNumber: requestPayload.SenderNumber,
				Receptor:     receptor,
				BatchId:      uuid.NewString(),
				LocalId:      requestPayload.LocalId,
				SendType:     1,
				SmsCount:     int(smsCount),
				Lang:         smsLanguage,
				Message:      requestPayload.Message,
				Status:       1,
				StatusText:   utils.StatusText(1),
				Date:         time.Now().Unix(),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			// Save the message in to database
			messageid, _, err := app.Models.SendSMS.Insert(SMSRow)
			if err != nil {
				log.Println("Can't be store Message in database")
				continue
			}
			// Add SMS to array for provider
			SMSRow.ID = messageid.InsertedID.(primitive.ObjectID)
			SMSArray = append(SMSArray, SMSRow)
			// Add the message info to response payload
			responsePayload = append(responsePayload, SendSmsResponse{
				Sender:       SMSRow.Sender,
				SenderNumber: SMSRow.SenderNumber,
				Receptor:     SMSRow.Receptor,
				SmsCount:     SMSRow.SmsCount,
				Lang:         SMSRow.Lang,
				Message:      SMSRow.Message,
				Status:       SMSRow.Status,
				StatusText:   SMSRow.StatusText,
				BatchId:      SMSRow.BatchId,
				LocalId:      SMSRow.LocalId,
				Date:         SMSRow.Date,
			})
		} else {
			log.Printf("mobile number is not valid ! mobile : %v", receptor)
		}
	}

	// Send message with oprator
	switch SMSArray[0].Sender {
	case "rahyab":
		provider.RahyabSendSms(SMSArray)
	case "rahyabPG":
		provider.PGSendSms(SMSArray)
	case "kavenegar":
		provider.KVSendSms(SMSArray)
	case "hamyarsms":
		provider.HamyarSMSSendSms(SMSArray)
	default:
		utils.ErrorJSON(w, errors.New("sender is not valid ! "), http.StatusBadRequest)
	}

	// Return Response for client

	log.Printf("All Messages sent successfully !")

	res := utils.JsonResponse{
		Error:   false,
		Message: "All Messages sent successfully !",
		Data:    responsePayload,
	}

	utils.WriteJSON(w, http.StatusAccepted, res)

}

func SendSMSArray(w http.ResponseWriter, r *http.Request) {
	var localId string
	mobileRegex, _ := regexp.Compile(`^(?:98|\+98|0098|0)?9[0-9]{9}$`)
	var SMSArray []data.SendSMS
	var responsePayload []SendSmsResponse
	var requestPayload struct {
		Message      []string `json:"message"`
		Receptor     []string `json:"receptor"`
		Sender       string   `json:"sender"`
		SenderNumber string   `json:"sender_number"`
		LocalId      []string `json:"localid,omitempty"`
	}

	// Check the request payload
	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check the request payload are not null
	if len(requestPayload.Receptor) == 0 || len(requestPayload.Message) == 0 || len(strings.TrimSpace(requestPayload.Sender)) == 0 || len(strings.TrimSpace(requestPayload.SenderNumber)) == 0 {
		utils.ErrorJSON(w, errors.New("invalid parameter ! parameters could not be null"), http.StatusBadRequest)
		return
	}

	// Check the message,receptor and localId array length are equal
	if len(requestPayload.Receptor) != len(requestPayload.Message) || (len(requestPayload.LocalId) > 0 && len(requestPayload.LocalId) != len(requestPayload.Message)) {
		utils.ErrorJSON(w, errors.New("invalid parameter ! array parameters length aren't equal"), http.StatusBadRequest)
		return
	}

	//  Check  the receptor array length are not more than 100
	if len(requestPayload.Receptor) > 100 {
		utils.ErrorJSON(w, errors.New("receptors array not be more than 100 number"), http.StatusBadRequest)
		return
	}

	for index, receptor := range requestPayload.Receptor {
		smsCount, smsLanguage, err := utils.SmsCount(requestPayload.Message[index])
		if err != nil {
			log.Println("message length is too high")
			continue
		}

		if mobileRegex.MatchString(receptor) {
			if len(requestPayload.LocalId) > 0 {
				localId = requestPayload.LocalId[index]
			} else {
				localId = ""
			}
			SMSRow := data.SendSMS{
				Sender:       requestPayload.Sender,
				SenderNumber: requestPayload.SenderNumber,
				Receptor:     receptor,
				BatchId:      uuid.NewString(),
				LocalId:      localId,
				SendType:     1,
				SmsCount:     int(smsCount),
				Lang:         smsLanguage,
				Message:      requestPayload.Message[index],
				Status:       1,
				StatusText:   utils.StatusText(1),
				Date:         time.Now().Unix(),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			// Save the message in to database
			id, _, err := app.Models.SendSMS.Insert(SMSRow)

			if err != nil {
				log.Println("Can't be store Message in database")
				continue
			}
			// Add SMS to array for provider
			SMSRow.ID = id.InsertedID.(primitive.ObjectID)
			SMSArray = append(SMSArray, SMSRow)

			// Add the message info to response payload
			responsePayload = append(responsePayload, SendSmsResponse{
				Sender:       SMSRow.Sender,
				SenderNumber: SMSRow.SenderNumber,
				Receptor:     SMSRow.Receptor,
				SmsCount:     SMSRow.SmsCount,
				Lang:         SMSRow.Lang,
				Message:      SMSRow.Message,
				Status:       SMSRow.Status,
				StatusText:   SMSRow.StatusText,
				BatchId:      SMSRow.BatchId,
				LocalId:      SMSRow.LocalId,
				Date:         SMSRow.Date,
			})
		} else {
			log.Printf("mobile number is not valid ! mobile : %v", receptor)
		}
	}

	// Send message with oprator
	switch SMSArray[0].Sender {
	case "rahyab":
		provider.RahyabSendSMSArray(SMSArray)
	case "rahyabPG":
		provider.PGSendSMSArray(SMSArray)
	case "kavenegar":
		provider.KVSendSMSArray(SMSArray)
	case "hamyarsms":
		provider.HamyarSMSSendSmsArray(SMSArray)
	default:
		utils.ErrorJSON(w, errors.New("sender is not valid ! "), http.StatusBadRequest)
	}

	// Return Response for client
	log.Printf("All Messages sent successfully !")

	res := utils.JsonResponse{
		Error:   false,
		Message: "All Messages sent successfully !",
		Data:    responsePayload,
	}

	utils.WriteJSON(w, http.StatusAccepted, res)

}
