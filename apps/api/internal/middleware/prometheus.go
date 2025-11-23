package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjain/pixelflow/apps/api/internal/metrics"
)

// PrometheusMiddleware records HTTP request metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		// Don't record metrics for the /metrics endpoint itself to avoid noise
		if path != "/metrics" {
			metrics.RequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
			metrics.RequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
		}
	}
}
