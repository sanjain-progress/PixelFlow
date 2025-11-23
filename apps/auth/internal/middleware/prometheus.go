package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjain/pixelflow/apps/auth/internal/metrics"
)

// PrometheusMiddleware records HTTP metrics for each request
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics after request completes
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		endpoint := c.FullPath()
		method := c.Request.Method

		// Record request count
		metrics.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()

		// Record request duration
		metrics.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}
