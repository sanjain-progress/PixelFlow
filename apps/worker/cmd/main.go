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
	"github.com/sanjain/pixelflow/apps/worker/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	// Initialize Tracer
	shutdown := tracing.InitTracer("worker-service")
	defer shutdown(context.Background())

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

		// Extract Trace Context
		carrier := propagation.MapCarrier{}
		for _, h := range event.Headers {
			carrier[h.Key] = string(h.Value)
		}
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

		// Start Span
		tracer := otel.Tracer("worker-service")
		_, span := tracer.Start(ctx, "process_task")
		defer span.End()
		
		// Increment consumed metric
		// Tracks total messages pulled from Kafka, regardless of processing outcome
		metrics.KafkaMessagesConsumedTotal.Inc()
		
		// Track active tasks using a Gauge
		// This helps us see if the worker is overwhelmed or stuck
		metrics.ActiveProcessingTasks.Inc()
		defer metrics.ActiveProcessingTasks.Dec()

		start := time.Now()

		// Process the image
		err := proc.ProcessImage(event.TaskID)
		
		// Record processing duration
		// We use a Histogram to calculate percentiles (P95, P99) later
		duration := time.Since(start).Seconds()
		metrics.TaskProcessingDuration.Observe(duration)

		if err != nil {
			slog.Error("Failed to process task", "task_id", event.TaskID, "error", err)
			// Track failures separately for error rate calculation
			metrics.TasksProcessedTotal.WithLabelValues("failure").Inc()
			metrics.KafkaConsumptionErrorsTotal.Inc()
			return err
		}

		slog.Info("Task completed successfully", "task_id", event.TaskID)
		metrics.TasksProcessedTotal.WithLabelValues("success").Inc()
		return nil
	})
}
