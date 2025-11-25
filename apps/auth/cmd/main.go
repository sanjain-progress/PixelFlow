package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sanjain/pixelflow/apps/auth/internal/db"
	"github.com/sanjain/pixelflow/apps/auth/internal/metrics"
	"github.com/sanjain/pixelflow/apps/auth/internal/middleware"
	"github.com/sanjain/pixelflow/apps/auth/internal/models"
	"github.com/sanjain/pixelflow/apps/auth/internal/tracing"
	"github.com/sanjain/pixelflow/apps/auth/internal/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	shutdown := tracing.InitTracer("auth-service")
	defer shutdown(context.Background())

	// Load Configuration
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/auth_db?sslmode=disable")
	port := getEnv("PORT", "50051")

	slog.Info("Starting Auth Service", "port", port)

	// Initialize Database Connection
	h := db.Init(dbURL)

	// Setup Gin HTTP server
	r := gin.Default()

	// Add OpenTelemetry Middleware
	r.Use(otelgin.Middleware("auth-service"))

	// Add Prometheus Middleware
	r.Use(middleware.PrometheusMiddleware())

	// CORS Middleware
	// Allows frontend (localhost:3000) to communicate with this service
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

	// Prometheus Metrics Endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// POST /register
	// Creates a new user account with hashed password
	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("Register: Invalid request body", "error", err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Check if user exists
		var existing models.User
		if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			slog.Warn("Register: User already exists", "email", req.Email)
			c.JSON(400, gin.H{"error": "User already exists"})
			return
		}

		// Hash password
		hashedPw, err := utils.HashPassword(req.Password)
		if err != nil {
			slog.Error("Register: Failed to hash password", "error", err)
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		// Create user
		user := models.User{
			Email:    req.Email,
			Password: hashedPw,
		}

		if err := h.DB.Create(&user).Error; err != nil {
			slog.Error("Register: Failed to create user in DB", "error", err)
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}

		// Record business metric
		metrics.RegistrationsTotal.Inc()

		slog.Info("User registered successfully", "email", req.Email, "user_id", user.ID)
		c.JSON(200, gin.H{"message": "User registered successfully"})
	})

	// POST /login
	// Authenticates user and returns JWT token
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("Login: Invalid request body", "error", err)
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Find user
		var user models.User
		if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			slog.Warn("Login: User not found", "email", req.Email)
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Check password
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			metrics.LoginsTotal.WithLabelValues("failure").Inc()
			slog.Warn("Login: Invalid password", "email", req.Email)
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT
		token, err := utils.GenerateJWT(fmt.Sprintf("%d", user.ID))
		if err != nil {
			slog.Error("Login: Failed to generate token", "error", err)
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		// Record successful login
		metrics.LoginsTotal.WithLabelValues("success").Inc()

		slog.Info("User logged in successfully", "email", req.Email, "user_id", user.ID)
		c.JSON(200, gin.H{"token": token})
	})

	// GET /validate
	// Validates JWT token and returns user ID
	r.GET("/validate", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			metrics.TokenValidationsTotal.WithLabelValues("invalid").Inc()
			c.JSON(401, gin.H{"valid": false})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			metrics.TokenValidationsTotal.WithLabelValues("invalid").Inc()
			c.JSON(401, gin.H{"valid": false})
			return
		}

		token := parts[1]
		userID, err := utils.ValidateJWT(token)
		if err != nil {
			metrics.TokenValidationsTotal.WithLabelValues("invalid").Inc()
			slog.Warn("Validate: Invalid token", "error", err)
			c.JSON(401, gin.H{"valid": false})
			return
		}

		// Record successful validation
		metrics.TokenValidationsTotal.WithLabelValues("valid").Inc()

		c.JSON(200, gin.H{"valid": true, "user_id": userID})
	})

	slog.Info("Auth Service listening", "address", ":"+port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
