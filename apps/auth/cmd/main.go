package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sanjain/pixelflow/apps/auth/internal/db"
	"github.com/sanjain/pixelflow/apps/auth/internal/models"
	"github.com/sanjain/pixelflow/apps/auth/internal/utils"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/auth_db?sslmode=disable")
	port := getEnv("PORT", "50051")

	h := db.Init(dbURL)

	// Setup Gin HTTP server
	r := gin.Default()

	// Register endpoint
	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Check if user exists
		var existing models.User
		if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			c.JSON(400, gin.H{"error": "User already exists"})
			return
		}

		// Hash password
		hashedPw, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		// Create user
		user := models.User{
			Email:    req.Email,
			Password: hashedPw,
		}

		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(200, gin.H{"message": "User registered successfully"})
	})

	// Login endpoint
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// Find user
		var user models.User
		if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Check password
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT
		token, err := utils.GenerateJWT(fmt.Sprintf("%d", user.ID))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	})

	// Validate endpoint
	r.GET("/validate", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"valid": false})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"valid": false})
			return
		}

		token := parts[1]
		userID, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(401, gin.H{"valid": false})
			return
		}

		c.JSON(200, gin.H{"valid": true, "user_id": userID})
	})

	fmt.Printf("Auth Service (HTTP) running on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
