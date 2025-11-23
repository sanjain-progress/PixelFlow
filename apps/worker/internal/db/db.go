package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handler wraps the MongoDB client and database
type Handler struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Init initializes MongoDB connection
func Init(mongoURL, dbName string) *Handler {
	ctx := context.Background()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	return &Handler{
		Client: client,
		DB:     client.Database(dbName),
	}
}
