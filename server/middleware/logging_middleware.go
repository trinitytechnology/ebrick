package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linkifysoft/ebrick/observability"
	"go.uber.org/zap"
)

// LoggingMiddleware logs requests with trace ID
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger := observability.LoggerWithTraceID(c.Request.Context())
		logger.Info("Request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}
