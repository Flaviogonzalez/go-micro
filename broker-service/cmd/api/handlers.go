package main

import (
	"net/http"
)

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
