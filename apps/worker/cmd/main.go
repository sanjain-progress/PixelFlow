package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sanjain/pixelflow/apps/worker/internal/db"
	"github.com/sanjain/pixelflow/apps/worker/internal/kafka"
	"github.com/sanjain/pixelflow/apps/worker/internal/metrics"
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
	metricsPort := getEnv("METRICS_PORT", "8081")

	slog.Info("Starting Worker Service", "kafka_brokers", kafkaBrokers)

	// 1. Start Metrics Server (Background)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		slog.Info("Metrics server listening", "port", metricsPort)
		if err := http.ListenAndServe(":"+metricsPort, nil); err != nil {
			slog.Error("Failed to start metrics server", "error", err)
		}
	}()

	// 2. Initialize MongoDB
	// Use the same DB as the API service
	dbHandler := db.Init(mongoURL, "pixelflow")
	slog.Info("Connected to MongoDB")

	// 3. Initialize Processor
	proc := processor.NewProcessor(dbHandler.DB)

	// 4. Initialize Kafka Consumer
	// GroupID "worker-group-1" ensures we can scale workers horizontally
	consumer := kafka.NewConsumer(
		kafkaBrokers,
		"image-tasks",
		"worker-group-1",
	)
	defer consumer.Close()
	slog.Info("Kafka Consumer initialized", "topic", "image-tasks", "group", "worker-group-1")

	// 5. Start Consuming
	// The handler function is called for each message
	slog.Info("Worker started consuming messages...")
	consumer.Consume(context.Background(), func(event kafka.TaskEvent) error {
		slog.Info("Received task", "task_id", event.TaskID, "user_id", event.UserID)
		
		// Increment consumed metric
		metrics.KafkaMessagesConsumedTotal.Inc()
		metrics.ActiveProcessingTasks.Inc()
		defer metrics.ActiveProcessingTasks.Dec()

		start := time.Now()

		// Process the image
		err := proc.ProcessImage(event.TaskID)
		duration := time.Since(start).Seconds()
		metrics.TaskProcessingDuration.Observe(duration)

		if err != nil {
			slog.Error("Failed to process task", "task_id", event.TaskID, "error", err)
			metrics.TasksProcessedTotal.WithLabelValues("failure").Inc()
			metrics.KafkaConsumptionErrorsTotal.Inc()
			return err
		}

		slog.Info("Task completed successfully", "task_id", event.TaskID)
		metrics.TasksProcessedTotal.WithLabelValues("success").Inc()
		return nil
	})
}
