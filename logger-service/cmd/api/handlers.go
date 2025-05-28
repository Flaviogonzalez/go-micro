package main

import (
	"log-service/cmd/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = app.ReadJSON(w, r, &requestPayload)

	// insert data

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.WriteJSON(w, http.StatusAccepted, res)
}
