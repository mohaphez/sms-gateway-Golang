package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (app *Config) getToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")
	// call the authentication service for check user credentials
	request, err := http.NewRequest("POST", "http://authentication-service:8081/verify-user", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"i": jsonFromService.Data.([]interface{})[0], "exp": time.Now().Add(time.Minute * 30)})

	if err != nil {
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
