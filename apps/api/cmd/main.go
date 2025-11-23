package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjain/pixelflow/apps/api/internal/db"
	"github.com/sanjain/pixelflow/apps/api/internal/kafka"
	"github.com/sanjain/pixelflow/apps/api/internal/middleware"
	"github.com/sanjain/pixelflow/apps/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "kafka:29092"), ",")
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:50051")
	port := getEnv("PORT", "8080")

	slog.Info("Starting API Service", "port", port, "kafka_brokers", kafkaBrokers)

	// 1. Initialize MongoDB
	// Connect to the 'pixelflow' database
	dbHandler := db.Init(mongoURL, "pixelflow")
	taskCollection := dbHandler.DB.Collection("tasks")
	slog.Info("Connected to MongoDB", "db", "pixelflow")

	// 2. Initialize Kafka Producer
	// Connect to Kafka broker and topic 'image-tasks'
	kafkaProducer := kafka.NewProducer(kafkaBrokers, "image-tasks")
	defer kafkaProducer.Close()
	slog.Info("Connected to Kafka Producer", "topic", "image-tasks")

	// 3. Initialize Auth Middleware
	// Connect to Auth Service gRPC server
	authMiddleware, err := middleware.NewAuthMiddleware(authServiceURL)
	if err != nil {
		slog.Error("Failed to connect to Auth Service", "error", err)
		os.Exit(1)
	}

	// 4. Setup Router
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Protected Routes (Require Authentication)
	// Apply auth middleware to protected routes
	authRoutes := r.Group("/api").Use(authMiddleware.Middleware())
	{
		// POST /api/upload - Create a new task
		authRoutes.POST("/upload", func(c *gin.Context) {
			// Get UserID from context (set by middleware)
			userID := c.GetString("userID")

			// Parse request body
			var req struct {
				ImageURL string `json:"image_url" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				slog.Warn("Upload: Invalid request", "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Create Task object
			task := models.Task{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				ImageURL:    req.ImageURL,
				Status:      models.StatusPending,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Save to MongoDB
			_, err := taskCollection.InsertOne(context.Background(), task)
			if err != nil {
				slog.Error("Upload: Failed to save task", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
				return
			}

			// Publish event to Kafka
			err = kafkaProducer.PublishTask(context.Background(), kafka.TaskEvent{
				TaskID:      task.ID.Hex(),
				UserID:      task.UserID,
				ImageURL:    task.ImageURL,
			})
			if err != nil {
				// Note: In production, we might want to rollback the DB insert or retry
				slog.Error("Upload: Failed to publish to Kafka", "error", err)
			} else {
				slog.Info("Task published to Kafka", "task_id", task.ID.Hex())
			}

			c.JSON(http.StatusCreated, task)
		})

		// GET /api/tasks - List user's tasks
		authRoutes.GET("/tasks", func(c *gin.Context) {
			userID := c.GetString("userID")

			// Find tasks for this user
			cursor, err := taskCollection.Find(context.Background(), bson.M{"user_id": userID})
			if err != nil {
				slog.Error("ListTasks: DB Query failed", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
				return
			}
			defer cursor.Close(context.Background())

			var tasks []models.Task
			if err = cursor.All(context.Background(), &tasks); err != nil {
				slog.Error("ListTasks: Decode failed", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tasks"})
				return
			}

			c.JSON(http.StatusOK, tasks)
		})
	}

	// 5. Start Server
	slog.Info("API Service listening", "address", ":"+port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
