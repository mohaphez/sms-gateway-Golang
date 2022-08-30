package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sms-service/cmd/api/utils"
	"sms-service/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Credential struct {
	Token string
}

var RahyabToken Credential

func RahyabGetToken() error {
	requestPayload := `{
		"Username": "",
		"Password": "",
		"Company":  ""
	}`
	url := "https://api.rahyab.ir/api/Auth/getToken"
	res, statusCode := utils.SendPostRequest(url, requestPayload, "")
	if statusCode == http.StatusOK {
		RahyabToken.Token = res
	}
	time.AfterFunc(24*time.Hour, func() { RahyabGetToken() })
	return nil
}

func RahyabSendSms(id *mongo.InsertOneResult, sms data.SendSMS) error {
	type responsePayload struct {
		Status       string
		Code         string
		Identity     string
		ErrorMsg     string
		ErrorMessage string
	}

	// Create request json for post request .
	requestPayload := fmt.Sprintf(`{
		"message": "%s",
	    "destinationAddress":"%s",
	    "number": "%s",
		"userName": "",
		"password": "",
		"company":  ""
	}`, sms.Message, sms.Receptor, sms.SenderNumber)

	// Check and fill token if it's empty
	if RahyabToken.Token == "" {
		RahyabGetToken()
	}

	// send sms for rahyab provider .

	url := "https://api.rahyab.ir/api/v1/SendSMS_Single"
	res, statusCode := utils.SendPostRequest(url, requestPayload, RahyabToken.Token)
	var responseJson responsePayload

	if statusCode == http.StatusOK {
		err := json.Unmarshal([]byte(res), &responseJson)
		if err != nil {
			panic(err)
		}

		// check provider response and update sms status .
		if len(responseJson.Status) == 0 {
			if len(responseJson.ErrorMsg) == 0 {
				sms.Identity = responseJson.Identity
				sms.Status = 3
			} else {
				sms.Status = 6
				sms.Error = responseJson.ErrorMsg
			}
		} else {
			sms.Status = 7
			sms.Error = responseJson.ErrorMessage
		}
		sms.StatusText = utils.StatusText(int16(sms.Status))
		sms.SendTime = time.Now()
		sms.UpdatedAt = time.Now()
		sms.Update(id, sms)
	}
	return nil
}
