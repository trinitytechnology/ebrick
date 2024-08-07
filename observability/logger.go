package observability

import (
	"context"

	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// LoggerWithTraceID returns a logger with trace ID from the context.
func LoggerWithTraceID(ctx context.Context) *zap.Logger {
	if config.GetConfig().Observability.Tracing.Enable {
		span := trace.SpanFromContext(ctx)
		if span == nil {
			return logger.DefaultLogger
		}
		sc := span.SpanContext()
		return logger.DefaultLogger.With(
			zap.String("trace_id", sc.TraceID().String()),
			zap.String("span_id", sc.SpanID().String()),
		)
	} else {
		return logger.DefaultLogger
	}
}
