package handler

import (
	"authentication-service/cmd/api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Authentication(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user exist in db
	user, password, err := app.Models.User.UserIsActive(requestPayload.UserName)

	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := app.Models.User.PasswordMatch(requestPayload.Password, password)

	if !valid || err != nil {
		utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	res := utils.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.UserName),
		Data:    user,
	}

	utils.WriteJSON(w, http.StatusAccepted, res)

}

func Createuser(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name,omitempty"`
		UserName  string    `json:"username,omitempty"`
		Password  string    `json:"password"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user exist in db
	user, err := app.Models.User.Insert(requestPayload)

	if err != nil {
		utils.ErrorJSON(w, errors.New("parameters are not valid "), http.StatusBadRequest)
		return
	}
	log.Printf("New user created successfully - user_id:%v", user)

	res := utils.JsonResponse{
		Error:   false,
		Message: "New user Created successfully",
	}

	utils.WriteJSON(w, http.StatusAccepted, res)

}
