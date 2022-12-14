package provider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sms-service/cmd/api/utils"
	"sms-service/data"
	"strconv"
	"strings"
	"time"
)

func HamyarSMSSendSoap(url string, method string, username string, password string, to []string, from string, msg string, flash bool, udh string, recid []int64) ([]byte, error) {

	_to := "<ns2:string>" + strings.Join(to, "</ns2:string><ns2:string>") + "</ns2:string>"
	_recid := "<ns2:long>" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(recid)), "</ns2:long><ns2:long>"), "[]") + "</ns2:long>"

	wsReq := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
			<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://tempuri.org/" xmlns:ns2="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
				<SOAP-ENV:Body>
					<ns1:%s>
						<ns1:userName>%s</ns1:userName>
						<ns1:password>%s</ns1:password>
						<ns1:fromNumber>%s</ns1:fromNumber>
						<ns1:toNumbers>%v</ns1:toNumbers>
						<ns1:messageContent>%s</ns1:messageContent>
						<ns1:isFlash>%s</ns1:isFlash>
						<ns1:recId>%v</ns1:recId>
						<ns1:status>MA==</ns1:status>
					</ns1:%s>
				</SOAP-ENV:Body>
			</SOAP-ENV:Envelope>
`, method, username, password, from, _to, msg, strconv.FormatBool(flash), _recid, method)

	dataAsByte := []byte(wsReq)
	req, err := http.NewRequest("Post", url, bytes.NewBuffer(dataAsByte))
	req.Header.Set("SOAPAction", fmt.Sprintf("http://tempuri.org/ISendService/%s", method))
	req.Header.Set("Content-Type", "text/xml")
	if err != nil {
		//Handle Error
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		log.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		fmt.Println(err)
		return []byte{52}, err
	}
	defer resp.Body.Close()
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		utils.LogEvent(logEvent)
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return []byte{0}, err
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		logEvent.Name = "info"
		logEvent.Data = fmt.Sprint(string(data))
		utils.LogEvent(logEvent)
		return data, nil
	}
}

func HamyarSMSSendSms(messages []data.SendSMS) error {

	// send sms for PGrahyab provider .
	for _, sms := range messages {
		res, err := HamyarSMSSendSoap(ProviderConfig.Providers.HamyarSMS.BaseUrl, "SendSMS", ProviderConfig.Providers.HamyarSMS.Username, ProviderConfig.Providers.HamyarSMS.Password, []string{sms.Receptor}, messages[0].SenderNumber, messages[0].Message, false, "", []int64{0})
		// Resend if request have connection error
		if err != nil && res[0] == 52 {
			logEvent.Name = "error"
			logEvent.Data = fmt.Sprint(err)
			utils.LogEvent(logEvent)
			time.Sleep(5 * time.Second)
			go HamyarSMSSendSms(messages)
			return nil
		}
		if err != nil {
			logEvent.Name = "error"
			logEvent.Data = fmt.Sprint(err)
			utils.LogEvent(logEvent)
			log.Println(err.Error())
			return err
		}
		// Check response status
		sendSmsResult, err := utils.UnmarshalXML(res, "SendSMSResult")
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
		if len(sendSmsResult) > 0 && sendSmsResult[0] == "0" {
			sms.Status = 3
			sms.Identity = recids[0]
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
			logEvent.Name = "info"
			logEvent.Data = fmt.Sprint(sendSmsResult)
		} else if len(sendSmsResult) > 0 {
			sms.Status = 11
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.Error = sendSmsResult[0]
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
			logEvent.Name = "warning"
			logEvent.Data = "provider error :" + fmt.Sprint(sendSmsResult)
			sms.Update(sms.ID, sms)
		}
		utils.LogEvent(logEvent)
	}
	return nil
}

func HamyarSMSSendSmsArray(messages []data.SendSMS) error {

	// send sms for PGrahyab provider .
	for _, sms := range messages {
		res, err := HamyarSMSSendSoap(ProviderConfig.Providers.HamyarSMS.BaseUrl, "SendSMS", ProviderConfig.Providers.HamyarSMS.Username, ProviderConfig.Providers.HamyarSMS.Password, []string{sms.Receptor}, messages[0].SenderNumber, sms.Message, false, "", []int64{0})
		// Resend if request have connection error
		if err != nil && res[0] == 52 {
			logEvent.Name = "error"
			logEvent.Data = fmt.Sprint(err)
			utils.LogEvent(logEvent)
			time.Sleep(5 * time.Second)
			go HamyarSMSSendSmsArray(messages)
			return nil
		}
		if err != nil {
			logEvent.Name = "error"
			logEvent.Data = fmt.Sprint(err)
			utils.LogEvent(logEvent)
			log.Println(err.Error())
			return err
		}
		// Check response status
		sendSmsResult, err := utils.UnmarshalXML(res, "SendSMSResult")
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
		if len(sendSmsResult) > 0 && sendSmsResult[0] == "0" {
			sms.Status = 3
			sms.Identity = recids[0]
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
			logEvent.Name = "info"
			logEvent.Data = fmt.Sprint(sendSmsResult)
		} else if len(sendSmsResult) > 0 {
			sms.Status = 11
			sms.StatusText = utils.StatusText(int16(sms.Status))
			sms.Error = sendSmsResult[0]
			sms.SendTime = time.Now()
			sms.UpdatedAt = time.Now()
			sms.Update(sms.ID, sms)
			logEvent.Name = "warning"
			logEvent.Data = "provider error :" + fmt.Sprint(sendSmsResult)
		}
		utils.LogEvent(logEvent)
	}
	return nil
}
