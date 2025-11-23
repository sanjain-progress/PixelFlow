package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjain/pixelflow/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthMiddleware validates JWT tokens using the Auth Service via gRPC.
type AuthMiddleware struct {
	authClient pb.AuthServiceClient
}

// NewAuthMiddleware creates a new middleware instance.
// authServiceUrl: Address of the Auth Service (e.g., "localhost:50051")
func NewAuthMiddleware(authServiceUrl string) (*AuthMiddleware, error) {
	// Connect to Auth Service
	// In production, use secure credentials (TLS)
	conn, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		authClient: pb.NewAuthServiceClient(conn),
	}, nil
}

// RequireAuth intercepts requests and checks for a valid Bearer token.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// 2. Extract token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		token := parts[1]

		// 3. Call Auth Service to validate token
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := m.authClient.Validate(ctx, &pb.ValidateRequest{Token: token})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 4. Set UserID in context for next handlers
		c.Set("userID", res.UserId)
		c.Next()
	}
}
