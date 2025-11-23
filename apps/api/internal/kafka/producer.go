package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer handles sending messages to Kafka.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer.
// brokers: List of Kafka broker addresses (e.g., ["localhost:9092"])
// topic: The topic to write to (e.g., "image-tasks")
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // Distribute messages evenly
	}

	fmt.Printf("Kafka Producer initialized for topic: %s\n", topic)
	return &Producer{writer: writer}
}

// TaskEvent represents the message sent to Kafka.
type TaskEvent struct {
	TaskID   string `json:"task_id"`
	UserID   string `json:"user_id"`
	ImageURL string `json:"image_url"`
}

// PublishTask sends a task event to Kafka.
func (p *Producer) PublishTask(ctx context.Context, event TaskEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Write message with a timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.TaskID), // Key ensures ordering for same task (if needed)
		Value: payload,
	})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	fmt.Printf("Published task event: %s\n", event.TaskID)
	return nil
}

// Close closes the producer connection.
func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Failed to close Kafka writer: %v", err)
	}
}
