package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handler manages the MongoDB connection.
type Handler struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Init establishes a connection to MongoDB.
// url: The MongoDB connection string (e.g., "mongodb://localhost:27017")
// dbName: The name of the database to use (e.g., "pixelflow")
func Init(url, dbName string) *Handler {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln("Failed to ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")

	return &Handler{
		Client: client,
		DB:     client.Database(dbName),
	}
}
