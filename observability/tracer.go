package observability

import (
	"context"

	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/zap"
)

func InitTracer(cfg config.TracingConfig) (*sdktrace.TracerProvider, error) {
	logger := logger.DefaultLogger
	logger.Info("Initializing tracer", zap.String("endpoint", cfg.Endpoint))
	serviceName := config.GetConfig().Service.Name
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(cfg.Endpoint))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
	logger.Info("Tracer initialized")
	return tp, nil
}
