package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	DB *mongo.Database
}

func Init(url, dbName string) *Handler {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalln(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to MongoDB")
	return &Handler{DB: client.Database(dbName)}
}
