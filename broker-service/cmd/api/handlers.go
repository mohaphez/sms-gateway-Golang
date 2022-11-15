package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		Error:   false,
		Message: "You Successfully Connected",
	}

	_ = app.writeJSON(w, http.StatusOK, res)
}

// Authentication user credentials
func (app *Config) GetToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	// call the authentication service for check user credentials
	request, err := http.NewRequest("POST", "http://authentication-service/verify-user", bytes.NewBuffer(jsonData))
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		logEvent.Name = "error"
		logEvent.Data = "invalid credentials"
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		logEvent.Name = "error"
		logEvent.Data = "error calling auth service"
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"i": jsonFromService.Data.([]interface{})[0], "exp": time.Now().Add(time.Minute * 30)})

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Println(tokenString)
	res := jsonResponse{
		Error:   false,
		Message: "You Successfully Authenticate",
		Data:    tokenString,
	}

	_ = app.writeJSON(w, http.StatusOK, res)
}

// Send SMS with sms service
func (app *Config) SendSMS(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Message      string   `json:"message"`
		Receptor     []string `json:"receptor"`
		Sender       string   `json:"sender"`
		SenderNumber string   `json:"sender_number"`
		LocalId      string   `json:"localid,omitempty"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	// call the sms service
	request, err := http.NewRequest("POST", "http://sms-service/send-sms", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, jsonFromService)
}

// Send P2P SMS with sms service
func (app *Config) SendSMSArray(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Message      []string `json:"message"`
		Receptor     []string `json:"receptor"`
		Sender       string   `json:"sender"`
		SenderNumber string   `json:"sender_number"`
		LocalId      []string `json:"localid,omitempty"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	// call the sms service
	request, err := http.NewRequest("POST", "http://sms-service/send-sms-p2p", bytes.NewBuffer(jsonData))
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, jsonFromService)
}

// Get SMS status from sms service
func (app *Config) SmsStatus(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Batchid []string `json:"batchid"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	// call the sms service
	request, err := http.NewRequest("POST", "http://sms-service/get-sms-status", bytes.NewBuffer(jsonData))
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(response.Body)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, jsonFromService)
}

// Get SMS list from sms service
func (app *Config) SmsList(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Message      string   `json:"message"`
		Receptor     []string `json:"receptor"`
		Sender       []string `json:"sender"`
		SenderNumber []string `json:"senderNumber"`
		Offset       int64    `json:"offset"`
		Limit        int64    `json:"limit"`
		Sort         string   `json:"sort"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	// call the sms service
	request, err := http.NewRequest("POST", "http://sms-service/get-sms-list", bytes.NewBuffer(jsonData))
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling sms service"))
		log.Println(response.Body)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		logEvent.Name = "error"
		logEvent.Data = fmt.Sprint(err)
		LogEvent(logEvent)
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, jsonFromService)
}
