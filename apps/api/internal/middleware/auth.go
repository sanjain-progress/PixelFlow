package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens by calling the Auth Service via HTTP
type AuthMiddleware struct {
	authServiceURL string
	httpClient     *http.Client
}

// ValidateResponse represents the response from auth service
type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(authServiceURL string) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		authServiceURL: authServiceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}, nil
}

// Close cleans up resources
func (m *AuthMiddleware) Close() error {
	return nil
}

// Middleware returns a Gin middleware handler that validates JWT tokens
func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Call Auth Service HTTP endpoint to validate token
		req, _ := http.NewRequest("GET", m.authServiceURL+"/validate", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := m.httpClient.Do(req)
		if err != nil {
			c.JSON(503, gin.H{"error": "Auth service unavailable"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var validateResp ValidateResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse auth response"})
			c.Abort()
			return
		}

		if !validateResp.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("userID", validateResp.UserID)
		c.Next()
	}
}
