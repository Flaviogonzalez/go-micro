package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{} // esto es como crear una caja y meter todas las rutas aca

func main() {
	// esto es toda la configuracion centrada en un lugar
	app := Config{}

	log.Printf("Startin broker service on port %s\n", webPort)

	// se define el server manualmente
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
