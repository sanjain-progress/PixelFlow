package processor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sanjain/pixelflow/apps/worker/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Processor handles the image processing logic.
type Processor struct {
	taskCollection *mongo.Collection
}

// NewProcessor creates a new processor instance.
func NewProcessor(db *mongo.Database) *Processor {
	return &Processor{
		taskCollection: db.Collection("tasks"),
	}
}

// ProcessImage simulates image processing and updates the task status.
func (p *Processor) ProcessImage(taskID string) error {
	ctx := context.Background()
	objID, _ := primitive.ObjectIDFromHex(taskID)

	// 1. Update Status to PROCESSING
	p.updateStatus(ctx, objID, models.StatusProcessing, "")
	fmt.Printf("Processing task: %s...\n", taskID)

	// 2. Simulate Processing (Sleep)
	time.Sleep(5 * time.Second)

	// 3. Generate "Processed" URL
	processedURL := fmt.Sprintf("https://cdn.pixelflow.com/processed/%s.jpg", taskID)

	// 4. Update Status to COMPLETED
	err := p.updateStatus(ctx, objID, models.StatusCompleted, processedURL)
	if err != nil {
		return err
	}

	fmt.Printf("Task %s completed!\n", taskID)
	return nil
}

// updateStatus helper to update MongoDB document.
func (p *Processor) updateStatus(ctx context.Context, id primitive.ObjectID, status models.TaskStatus, url string) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	// If URL is provided, update it too
	if url != "" {
		update["$set"].(bson.M)["processed_url"] = url
	}

	_, err := p.taskCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("Failed to update status for %s: %v", id.Hex(), err)
		return err
	}
	return nil
}
