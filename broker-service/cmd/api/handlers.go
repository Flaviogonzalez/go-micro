package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

/*
ok esto basicamente es una respuesta de servidor preparada manualmente.
completamente agradecido con el de arriba
*/
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request to / from %s", r.Method, r.RemoteAddr)
	payload := JsonResponse{
		Error:   false,
		Message: "Hit The Broker",
	}
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, `{"error": true, "message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write(out); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
