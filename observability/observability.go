package observability

import (
	"github.com/trinitytechnology/ebrick/config"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var DefaultTraceProvider *trace.TracerProvider = NewTracer()

// NewTracer creates a new instance of the tracer provider.
// It initializes the tracer based on the configuration settings and returns the tracer provider.
// If tracing is disabled in the configuration, it returns nil.
func NewTracer() *trace.TracerProvider {
	logger := zap.NewExample()
	cfg := config.GetConfig().Observability.Tracing

	if cfg.Enable {
		tp, err := InitTracer(cfg)
		if err != nil {
			logger.Error("Failed to initialize tracer", zap.Error(err))
		}
		return tp
	}

	return nil
}
