package main

import "net/http"

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		Error:   false,
		Message: "You Successfully Connected",
	}

	_ = app.writeJSON(w, http.StatusOK, res)
}
