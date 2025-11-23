package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjain/pixelflow/apps/api/internal/db"
	"github.com/sanjain/pixelflow/apps/api/internal/kafka"
	"github.com/sanjain/pixelflow/apps/api/internal/middleware"
	"github.com/sanjain/pixelflow/apps/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// 1. Initialize MongoDB
	// Connect to the 'pixelflow' database
	dbHandler := db.Init("mongodb://localhost:27017", "pixelflow")
	taskCollection := dbHandler.DB.Collection("tasks")

	// 2. Initialize Kafka Producer
	// Connect to Kafka broker and topic 'image-tasks'
	kafkaProducer := kafka.NewProducer([]string{"localhost:9092"}, "image-tasks")
	defer kafkaProducer.Close()

	// 3. Initialize Auth Middleware
	// Connect to Auth Service gRPC server
	authMiddleware, err := middleware.NewAuthMiddleware("localhost:50051")
	if err != nil {
		log.Fatalln("Failed to connect to Auth Service:", err)
	}

	// 4. Setup Router
	r := gin.Default()

	// Public Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Protected Routes (Require Authentication)
	api := r.Group("/api")
	api.Use(authMiddleware.RequireAuth())
	{
		// POST /api/upload - Create a new task
		api.POST("/upload", func(c *gin.Context) {
			// Get UserID from context (set by middleware)
			userID := c.GetString("userID")

			// Parse request body
			var req struct {
				ImageURL string `json:"image_url" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Create Task object
			task := models.Task{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				OriginalURL: req.ImageURL,
				Status:      models.StatusPending,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Save to MongoDB
			_, err := taskCollection.InsertOne(context.Background(), task)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
				return
			}

			// Publish event to Kafka
			err = kafkaProducer.PublishTask(context.Background(), kafka.TaskEvent{
				TaskID:      task.ID.Hex(),
				UserID:      task.UserID,
				OriginalURL: task.OriginalURL,
			})
			if err != nil {
				// Note: In production, we might want to rollback the DB insert or retry
				log.Println("Failed to publish to Kafka:", err)
			}

			c.JSON(http.StatusCreated, task)
		})

		// GET /api/tasks - List user's tasks
		api.GET("/tasks", func(c *gin.Context) {
			userID := c.GetString("userID")

			// Find tasks for this user
			cursor, err := taskCollection.Find(context.Background(), bson.M{"user_id": userID})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
				return
			}
			defer cursor.Close(context.Background())

			var tasks []models.Task
			if err = cursor.All(context.Background(), &tasks); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tasks"})
				return
			}

			c.JSON(http.StatusOK, tasks)
		})
	}

	// 5. Start Server
	fmt.Println("API Service running on :8080")
	r.Run(":8080")
}
