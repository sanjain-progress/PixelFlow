package main

import (
	"context"
	"fmt"
	"os"
	"strings"


	"github.com/sanjain/pixelflow/apps/worker/internal/db"
	"github.com/sanjain/pixelflow/apps/worker/internal/kafka"
	"github.com/sanjain/pixelflow/apps/worker/internal/processor"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Configuration
	mongoURL := getEnv("MONGO_URL", "mongodb://localhost:27017")
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9093"), ",")

	// 1. Initialize MongoDB
	// Use the same DB as the API service
	dbHandler := db.Init(mongoURL, "pixelflow")

	// 2. Initialize Processor
	proc := processor.NewProcessor(dbHandler.DB)

	// 3. Initialize Kafka Consumer
	// GroupID "worker-group-1" ensures we can scale workers horizontally
	consumer := kafka.NewConsumer(
		kafkaBrokers,
		"image-tasks",
		"worker-group-1",
	)
	defer consumer.Close()

	// 4. Start Consuming
	// The handler function is called for each message
	consumer.Consume(context.Background(), func(event kafka.TaskEvent) error {
		fmt.Printf("Worker received task: %s\n", event.TaskID)
		
		// Process the image
		return proc.ProcessImage(event.TaskID)
	})
}
