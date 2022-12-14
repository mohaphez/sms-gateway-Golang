package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sms-service/cmd/api/utils"
	"sms-service/data"
	"strings"
	"time"
)

type Credential struct {
	Token    string
	Username string
	Password string
	Company  string
}

var RahyabCredential = Credential{
	Token:    "",
	Username: "",
	Password: "",
	Company:  "",
}

func RahyabGetToken() error {
	requestPayload := fmt.Sprintf(`{
		"Username": "%s@%s",
		"Password": "%s",
		"Company":  "%s"
	}`, ProviderConfig.Providers.Rahyab.Username, ProviderConfig.Providers.Rahyab.Company, ProviderConfig.Providers.Rahyab.Password, ProviderConfig.Providers.Rahyab.Company)
	url := ProviderConfig.Providers.Rahyab.BaseUrl + "/Auth/getToken"
	res, statusCode := utils.SendPostRequest(url, requestPayload, "")
	if statusCode == http.StatusOK {
		RahyabCredential.Token = res
		RahyabCredential.Username = ProviderConfig.Providers.Rahyab.Username
		RahyabCredential.Password = ProviderConfig.Providers.Rahyab.Password
		RahyabCredential.Company = ProviderConfig.Providers.Rahyab.Company
	} else {
		logEvent.Name = "error"
		logEvent.Data = "Provider Status code :" + fmt.Sprint(statusCode) + " Provider Error:" + fmt.Sprint(res)
		utils.LogEvent(logEvent)
	}
	time.AfterFunc(24*time.Hour, func() { RahyabGetToken() })
	return nil
}

func RahyabSendSms(messages []data.SendSMS) error {
	type responseItem struct {
		SourceAddress string
		DestAddress   string
		Status        string
		Response      string
		SmsId         string
	}
	type responsePayload struct {
		SubmitResponse []responseItem
		BatchId        string
	}
	var destNumbersList []string
	for _, sms := range messages {
		destNumbersList = append(destNumbersList, fmt.Sprintf(`"%s"`, sms.Receptor))
	}
	// Create request json for post request .
	requestPayload := fmt.Sprintf(`{
		"message": "%s",
	    "destinationAddress":[%v],
	    "number": "%s",
		"userName": "%s",
		"password": "%s",
		"company":  "%s"
	}`, messages[0].Message, strings.Join(destNumbersList, ","), messages[0].SenderNumber, RahyabCredential.Username, RahyabCredential.Password, RahyabCredential.Company)

	// Check and fill token if it's empty
	if RahyabCredential.Token == "" {
		RahyabGetToken()
	}

	// send sms for rahyab provider .

	url := ProviderConfig.Providers.Rahyab.BaseUrl + "v1/SendSMS_Batch"
	res, statusCode := utils.SendPostRequest(url, requestPayload, RahyabCredential.Token)
	var responseJson []responsePayload
	// Resend if request have connection error
	if res == "" && statusCode == 502 {
		time.Sleep(5 * time.Second)
		go RahyabSendSms(messages)
		return nil
	}
	// Check response status
	if statusCode == http.StatusOK {
		err := json.Unmarshal([]byte(res), &responseJson)
		if err != nil {
			logEvent.Name = "error"
			logEvent.Data = fmt.Sprint(err)
			utils.LogEvent(logEvent)
			panic(err)
		}

		// check provider response and update sms status .
		for index, sms := range responseJson[0].SubmitResponse {
			if sms.Status == "CHECK_OK" {
				messages[index].Status = 3
			} else {
				messages[index].Status = 6
				messages[index].Error = sms.Response
			}

			messages[index].Identity = responseJson[0].BatchId
			messages[index].StatusText = utils.StatusText(int16(messages[index].Status))
			messages[index].SendTime = time.Now()
			messages[index].UpdatedAt = time.Now()
			messages[index].Update(messages[index].ID, messages[index])
		}
	}
	return nil
}

func RahyabSendSMSArray(messages []data.SendSMS) error {
	type responsePayload struct {
		Status       string
		Code         string
		Identity     string
		ErrorMsg     string
		ErrorMessage string
	}

	var LikeToLikeMessageList []string

	for _, sms := range messages {
		message := fmt.Sprintf(`{
			"message": "%s",
	    	"destNumber": "%s",
			"messageId": "%s"
		}`, sms.Message, sms.Receptor, sms.LocalId)

		LikeToLikeMessageList = append(LikeToLikeMessageList, message)
	}

	// Create request json for post request .
	requestPayload := fmt.Sprintf(`{
		"listLikeToLikeMessage": [%v],
	    "number": "%s",
		"userName": "%s",
		"password": "%s",
		"company":  "%s"
	}`, strings.Join(LikeToLikeMessageList, ","), messages[0].SenderNumber, RahyabCredential.Username, RahyabCredential.Password, RahyabCredential.Company)

	// Check and fill token if it's empty
	if RahyabCredential.Token == "" {
		RahyabGetToken()
	}

	// send sms for rahyab provider .

	url := ProviderConfig.Providers.Rahyab.BaseUrl + "v1/SendSMS_LikeToLike"
	res, statusCode := utils.SendPostRequest(url, requestPayload, RahyabCredential.Token)
	var responseJson []responsePayload
	// Resend if request have connection error
	if res == "" && statusCode == 502 {
		time.Sleep(5 * time.Second)
		go RahyabSendSMSArray(messages)
		return nil
	}

	if statusCode == http.StatusOK {
		err := json.Unmarshal([]byte(res), &responseJson)
		if err != nil {
			logEvent.Name = "error"
			logEvent.Data = fmt.Sprint(err)
			utils.LogEvent(logEvent)
			log.Println(err)
			return nil
		}

		for _, sms := range messages {
			// check provider response and update sms status .
			if len(responseJson[0].Status) == 0 {
				if len(responseJson[0].ErrorMsg) == 0 {
					sms.Identity = responseJson[0].Identity
					sms.Status = 3
				} else {
					sms.Status = 6
					sms.Error = responseJson[0].ErrorMsg
				}
			} else {
				sms.Status = 7
				sms.Error = responseJson[0].ErrorMessage
			}
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
		}

	}
	return nil
}
