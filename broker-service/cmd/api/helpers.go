package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // un mega

	// el limite de lectura para un json es un miserable mega
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// convierte de bit64 a JSON
	dec := json.NewDecoder(r.Body)

	// decodifica de JSON a golang structs
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	/*
		1 - &struct{}{} es el ANY de los structs
		2 - dec.Decode viene de decodificar un json pero espera que ese sea el único valor.
		3 - si hay mas de un valor en el request, se rellena en &struct{}{} || ejemplo : cae una request application/json de esta manera:
		*** {name: pirulo, password: asd} <- esta perfecto, es un solo JSON
		*** {name: pirulo, password: asd}{asd: pirulo, virus: implemented{}} <- esta mal, contenido malicioso detectado
		4 - entonces el error deberia ser que el "end of input" (eof) termina en el JSON y no en los datos que no deberían pasar
	*/
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("te estas pasando de vivo maestro")
	}

	return nil
}

func (app *Config) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	return app.WriteJSON(w, statusCode, payload)

}
