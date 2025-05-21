package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webport = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service...")

	// 1. conectar a la db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Mensaje desde el main: no se puede conectar a postgres todavia")
	}

	// 2. empezar la config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("postgres no parece estar listo...")
			counts++
		} else {
			log.Println("conectado a postgress")

			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("esperando por dos segundos")
		time.Sleep(2 * time.Second)
		continue
	}
}
