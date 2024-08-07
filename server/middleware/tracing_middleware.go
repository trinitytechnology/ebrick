package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/linkifysoft/ebrick/config"
	"go.opentelemetry.io/otel"
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
