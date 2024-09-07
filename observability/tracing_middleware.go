package observability

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linkifysoft/ebrick/config"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// TracingMiddleware adds tracing to the request context
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := config.GetConfig().Service.Name
		tracer := otel.Tracer(serviceName)
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func LoggingWithTraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger := LoggerWithTraceID(c.Request.Context())
		logger.Info("Request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}
