package main

import (
	"context"
	"fmt"
	"log"
	"log-service/cmd/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting logger-service...")

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic("Failed to connect to MongoDB:", err)
	}
	client = mongoClient
	log.Println("Connected to MongoDB")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	defer func() {
		log.Println("Disconnecting from MongoDB...")
		if err = client.Disconnect(ctx); err != nil {
			log.Panic("Error disconnecting MongoDB:", err)
		}
		log.Println("Disconnected from MongoDB")
	}()

	app := Config{
		Models: data.New(client),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting HTTP server on port %s\n", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("HTTP server error:", err)
	}
}

func (app Config) serve() {
	log.Printf("Starting HTTP server on port %s (serve method)\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic("HTTP server error (serve method):", err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	log.Println("Connecting to MongoDB...")
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v\n", err)
		return nil, err
	}

	log.Println("MongoDB connection established")
	return conn, nil
}
