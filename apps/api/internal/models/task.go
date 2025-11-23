package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskStatus defines the possible states of a task.
type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusProcessing TaskStatus = "PROCESSING"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

// Task represents an image processing job.
// It is stored in MongoDB and tracks the lifecycle of an upload.
type Task struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"user_id" json:"user_id"`             // ID of the user who uploaded the image
	OriginalURL  string             `bson:"original_url" json:"original_url"`   // URL of the raw uploaded image
	ProcessedURL string             `bson:"processed_url" json:"processed_url"` // URL of the processed image (empty initially)
	Status       TaskStatus         `bson:"status" json:"status"`               // Current status (PENDING, COMPLETED, etc.)
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
