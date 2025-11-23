package main

import (
	"context"
	"log/slog"
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
	// Initialize Structured Logger (JSON)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Configuration
	mongoURL := getEnv("MONGO_URL", "mongodb://localhost:27017")
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9093"), ",")

	slog.Info("Starting Worker Service", "kafka_brokers", kafkaBrokers)

	// 1. Initialize MongoDB
	// Use the same DB as the API service
	dbHandler := db.Init(mongoURL, "pixelflow")
	slog.Info("Connected to MongoDB")

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
	slog.Info("Kafka Consumer initialized", "topic", "image-tasks", "group", "worker-group-1")

	// 4. Start Consuming
	// The handler function is called for each message
	slog.Info("Worker started consuming messages...")
	consumer.Consume(context.Background(), func(event kafka.TaskEvent) error {
		slog.Info("Received task", "task_id", event.TaskID, "user_id", event.UserID)

		// Process the image
		err := proc.ProcessImage(event.TaskID)
		if err != nil {
			slog.Error("Failed to process task", "task_id", event.TaskID, "error", err)
			return err
		}

		slog.Info("Task completed successfully", "task_id", event.TaskID)
		return nil
	})
}
