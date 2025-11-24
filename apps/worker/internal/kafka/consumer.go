package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// TaskEvent represents the message received from Kafka.
type TaskEvent struct {
	TaskID      string         `json:"task_id"`
	UserID      string         `json:"user_id"`
	OriginalURL string         `json:"original_url"`
	Headers     []kafka.Header `json:"-"`
}

// Consumer handles reading messages from Kafka.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer.
// brokers: List of Kafka broker addresses
// topic: Topic to consume from
// groupID: Consumer group ID (for load balancing)
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
		GroupID:     groupID,
	})

	fmt.Printf("Kafka Consumer initialized for topic: %s (Group: %s)\n", topic, groupID)
	return &Consumer{reader: reader}
}

// Consume starts the consumer loop.
// handler: A function that processes each received task.
//
// IMPORTANT: This function implements retry logic to handle Kafka connection failures.
// Common scenario: Worker starts before Kafka is fully ready during docker-compose startup.
// Instead of crashing on first error, we log and continue trying to read messages.
func (c *Consumer) Consume(ctx context.Context, handler func(TaskEvent) error) {
	fmt.Println("Worker started consuming messages...")

	for {
		// 1. Read Message
		// Note: ReadMessage blocks until a message is available or an error occurs
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			// CHANGED: Instead of breaking (which exits the loop), we log and continue
			// This handles temporary Kafka connection issues during startup
			// Example: Worker starts at T+0s, Kafka ready at T+6s
			// Without retry: Worker crashes at T+3s when first connection fails
			// With retry: Worker keeps trying and succeeds when Kafka is ready
			slog.Info("Failed to read message: " + err.Error())
			continue // Keep trying instead of breaking
		}

		// 2. Parse Message
		var event TaskEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			slog.Warn("Failed to unmarshal event", "error", err)
			continue // Skip malformed messages
		}
		event.Headers = m.Headers

		fmt.Printf("Received task: %s\n", event.TaskID)

		// 3. Process Message (Call the handler)
		if err := handler(event); err != nil {
			slog.Error("Failed to process task", "task_id", event.TaskID, "error", err)
			// Note: In a real app, we might want to retry or send to a Dead Letter Queue (DLQ)
		}
	}
}

// Close closes the consumer connection.
func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Printf("Failed to close Kafka reader: %v", err)
	}
}
