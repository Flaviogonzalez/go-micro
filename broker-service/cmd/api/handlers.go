package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
ok esto basicamente es una respuesta de servidor preparada manualmente.
completamente agradecido con el de arriba
*/
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit The Broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload

	err := app.ReadJSON(w, r, &requestPayload) // nunca olvidar el &&&&&&&&&&&&&&&&
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.ErrorJSON(w, errors.New("unknown action"))
		return
	}
}

func (app *Config) authenticate(w http.ResponseWriter, auth AuthPayload) {
	log.Println("Starting authentication process")

	// 1. crear un struct json que voy a mandar al authentication service
	jsonData, err := json.MarshalIndent(auth, "", "\t")
	if err != nil {
		log.Printf("Error marshaling auth payload: %v", err)
		app.ErrorJSON(w, err)
		return
	}
	log.Printf("Auth payload marshaled: %s", string(jsonData))

	// 2. llamar al servicio
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		app.ErrorJSON(w, err)
		return
	}
	log.Println("HTTP request created for authentication service")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()
	log.Printf("Received response from authentication service: %d", response.StatusCode)

	// 3. y devolver el codigo de estado correcto
	if response.StatusCode == http.StatusUnauthorized {
		log.Println("Unauthorized credentials")
		app.ErrorJSON(w, errors.New("unauthorized credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Printf("Unexpected status code: %d", response.StatusCode)
		app.ErrorJSON(w, errors.New("something went wrong"))
		return
	}

	// 4. crear una variable para meterle response.body dentro
	var jsonFromService jsonResponse

	// 5. decode
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Printf("Error decoding response body: %v", err)
		app.ErrorJSON(w, err)
		return
	}
	log.Printf("Decoded response from service: %+v", jsonFromService)

	if jsonFromService.Error {
		log.Println("Authentication service returned error")
		app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	log.Println("Authentication successful, sending response")
	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}
