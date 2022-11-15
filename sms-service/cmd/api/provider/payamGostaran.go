package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"sms-service/cmd/api/utils"
	"sms-service/data"
	"time"
)

// Send batch Messages with Rahyab Payam Gostaran  provider

func PGSendSms(messages []data.SendSMS) error {

	var destNumbersList []string
	for _, sms := range messages {
		destNumbersList = append(destNumbersList, sms.Receptor)
	}

	// send sms for PGrahyab provider .

	res, err := utils.SendSMSSoap(ProviderConfig.Providers.PG.BaseUrl, "SendSms", ProviderConfig.Providers.PG.Username, ProviderConfig.Providers.PG.Password, destNumbersList, messages[0].SenderNumber, messages[0].Message, false, "", []int64{0})
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		log.Println(err.Error())
		return err
	}

	// Check response status
	sendSmsResult, err := utils.UnmarshalXML(res, "SendSmsResult")
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		return err
	}
	recids, err := utils.UnmarshalXML(res, "long")
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		return err
	}
	// check provider response and update sms status .
	if len(sendSmsResult) > 0 && sendSmsResult[0] == "1" {
		for index, recId := range recids {
			messages[index].Status = 3
			messages[index].Identity = recId
			messages[index].StatusText = utils.StatusText(int16(messages[index].Status))
			messages[index].SendTime = time.Now()
			messages[index].UpdatedAt = time.Now()
			messages[index].Update(messages[index].ID, messages[index])
		}
	} else {
		for _, message := range messages {
			message.Status = 11
			message.StatusText = utils.StatusText(int16(message.Status))
			message.Error = sendSmsResult[0]
			message.SendTime = time.Now()
			message.UpdatedAt = time.Now()
			message.Update(message.ID, message)
		}
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(sendSmsResult)
		utils.LogEvent(logEvent)
	}

	return nil
}

// Send Like to Like Messages with Rahyab Payam Gostaran  provider
func PGSendSMSArray(messages []data.SendSMS) error {
	type responseJsonPayload struct {
		RecIDs  string `json:"recIDs,omitempty"`
		Success int    `json:"success,omitempty"`
	}

	var destNumbersList []string
	var destMessagesList []string
	for _, sms := range messages {
		destNumbersList = append(destNumbersList, sms.Receptor)
		destMessagesList = append(destMessagesList, sms.Message)
	}

	// send sms for PGrahyab provider .

	res, err := utils.SendSMSArraySoap(ProviderConfig.Providers.PG.BaseUrl, "CorrespondSMS", ProviderConfig.Providers.PG.Username, ProviderConfig.Providers.HamyarSMS.Password, destNumbersList, messages[0].SenderNumber, destMessagesList, false, "", []int64{0})
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		log.Println(err.Error())
		return err
	}
	log.Println(string(res))
	// Check response status
	sendSmsResult, err := utils.UnmarshalXML(res, "CorrespondSMSResult")
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		return err
	}
	var respnseJson responseJsonPayload
	var recIds []string
	json.Unmarshal([]byte(sendSmsResult[0]), &respnseJson)
	json.Unmarshal([]byte(respnseJson.RecIDs), &recIds)

	// check provider response and update sms status .
	for index, message := range messages {
		if len(recIds) > index {
			message.Identity = recIds[index]
			message.Status = 3
		} else {
			message.Error = "Error! short number is false"
			message.Status = 11
		}
		message.StatusText = utils.StatusText(int16(message.Status))
		message.SendTime = time.Now()
		message.UpdatedAt = time.Now()
		message.Update(message.ID, message)
	}

	return nil
}
