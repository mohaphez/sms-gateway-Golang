package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sms-service/cmd/api/utils"
	"sms-service/data"
	"strings"
	"time"
)

type KVResponsePayload struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []struct {
		Messageid  int64  `json:"messageid"`
		Message    string `json:"message"`
		Status     int16  `json:"status"`
		Statustext string `json:"statustext"`
		Sender     string `json:"sender"`
		Receptor   string `json:"receptor"`
		Date       int32  `json:"date"`
		Cost       int16  `json:"cost"`
	} `json:"entries,omitempty"`
}

func KVSendSms(messages []data.SendSMS) error {

	var destNumbersList []string
	for _, sms := range messages {
		destNumbersList = append(destNumbersList, sms.Receptor)
	}
	// Create request json for post request .
	var requestPayload = make(map[string]string)
	requestPayload["message"] = messages[0].Message
	requestPayload["sender"] = messages[0].SenderNumber
	requestPayload["receptor"] = strings.Join(destNumbersList, ",")

	// send sms for kevenegar provider .

	res, statusCode := utils.SendGetRequest(ProviderConfig.Providers.Kavenegar.BaseUrl+ProviderConfig.Providers.Kavenegar.ApiKey+"/sms/"+"send.json", requestPayload, "")
	var responseJson KVResponsePayload

	// Check response status
	err := json.Unmarshal([]byte(res), &responseJson)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		panic(err)
	}
	if statusCode == http.StatusOK {
		if responseJson.Return.Status == 200 {
			// check provider response and update sms status .
			for index, sms := range responseJson.Entries {
				messages[index].Identity = fmt.Sprint(sms.Messageid)
				messages[index].StatusText = utils.StatusText(int16(messages[index].Status))
				messages[index].SendTime = time.Now()
				messages[index].UpdatedAt = time.Now()
				messages[index].Update(messages[index].ID, messages[index])
			}
		}
	} else {
		for _, sms := range messages {
			sms.Status = 4
			sms.Error = responseJson.Return.Message
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
		}
		logEvent.Name = "warning"
		logEvent.Data = "provider error :" + fmt.Sprint(responseJson)
		utils.LogEvent(logEvent)
	}
	return nil
}

func KVSendSMSArray(messages []data.SendSMS) error {

	var messagesList []string
	var destNumbersList []string
	var sendNumbersList []string
	for _, sms := range messages {
		destNumbersList = append(destNumbersList, fmt.Sprintf(`"%s"`, sms.Receptor))
		sendNumbersList = append(sendNumbersList, fmt.Sprintf(`"%s"`, sms.SenderNumber))
		messagesList = append(messagesList, fmt.Sprintf(`"%s"`, sms.Message))
	}

	// Create request json for post request .
	var requestPayload = make(map[string]string)
	requestPayload["message"] = fmt.Sprintf(`[%v]`, strings.Join(messagesList, ","))
	requestPayload["sender"] = fmt.Sprintf(`[%v]`, strings.Join(sendNumbersList, ","))
	requestPayload["receptor"] = fmt.Sprintf(`[%v]`, strings.Join(destNumbersList, ","))
	// send sms for rahyab provider .

	res, statusCode := utils.SendPostFormRequest(ProviderConfig.Providers.Kavenegar.BaseUrl+ProviderConfig.Providers.Kavenegar.ApiKey+"/sms/"+"sendarray.json", requestPayload, "")
	var responseJson KVResponsePayload

	err := json.Unmarshal([]byte(res), &responseJson)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		panic(err)
	}
	if statusCode == http.StatusOK {
		if responseJson.Return.Status == 200 {
			// check provider response and update sms status .
			for index, sms := range responseJson.Entries {
				messages[index].Identity = fmt.Sprint(sms.Messageid)
				messages[index].StatusText = utils.StatusText(int16(messages[index].Status))
				messages[index].SendTime = time.Now()
				messages[index].UpdatedAt = time.Now()
				messages[index].Update(messages[index].ID, messages[index])
			}
		}
	} else {
		for _, sms := range messages {
			sms.Status = 4
			sms.Error = responseJson.Return.Message
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
		}
		logEvent.Name = "warning"
		logEvent.Data = "provider error :" + fmt.Sprint(responseJson)
		utils.LogEvent(logEvent)
	}
	return nil
}
