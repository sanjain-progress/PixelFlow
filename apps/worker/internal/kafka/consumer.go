package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// TaskEvent represents the message received from Kafka.
type TaskEvent struct {
	TaskID      string `json:"task_id"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
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
		Brokers:  brokers,
		Topic:    topic,
		StartOffset: kafka.FirstOffset,
	})

	fmt.Printf("Kafka Consumer initialized for topic: %s (Group: %s)\n", topic, groupID)
	return &Consumer{reader: reader}
}

// Consume starts the consumer loop.
// handler: A function that processes each received task.
func (c *Consumer) Consume(ctx context.Context, handler func(TaskEvent) error) {
	fmt.Println("Worker started consuming messages...")

	for {
		// 1. Read Message
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			break // Exit loop on error (or handle retry)
		}

		// 2. Parse Message
		var event TaskEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue // Skip malformed messages
		}

		fmt.Printf("Received task: %s\n", event.TaskID)

		// 3. Process Message (Call the handler)
		if err := handler(event); err != nil {
			log.Printf("Failed to process task %s: %v", event.TaskID, err)
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
