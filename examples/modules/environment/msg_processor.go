package environment

import (
	"context"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/linkifysoft/ebrick/entity"
	"github.com/linkifysoft/ebrick/messaging"
	"github.com/linkifysoft/ebrick/observability"
	"github.com/linkifysoft/ebrick/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Tenant struct {
	entity.AuditEntity
	Name string `json:"name"`
}

// configureConsumers sets up the consumers for the CloudEventStream to process incoming events.
func configureConsumers(stream messaging.CloudEventStream) {

	serviceName := "environment"
	// Start a new goroutine to handle the subscription
	if err := stream.CreateConsumerGroup("tenant", serviceName, messaging.ConsumerConfig{
		AckWait:        30 * time.Second,
		MaxDeliver:     5, // Maximum retry attempts
		DeliverGroup:   serviceName,
		DeliverSubject: "dlq.tenant", // Move to DLQ after max retries
	}); err != nil {
		log.Error("Failed to create consumer group", zap.Error(err))
	}

	go func() {
		err := stream.Subscribe("tenant.created", serviceName, func(ev *event.Event, ctx context.Context) error {
			// Logger with trace ID from the context
			log := observability.LoggerWithTraceID(ctx)

			// Extract the tracing context and start a new span
			tracer := otel.Tracer(serviceName)
			_, span := tracer.Start(ctx, "Process Tenant Created Event")
			span.SetAttributes(attribute.String("subject", ev.Type()), attribute.String("source", ev.Source()), attribute.String("module", "environment"))
			defer span.End()

			tent, err := utils.UnmarshalJSONByte[Tenant](ev.Data())

			if err != nil {
				log.Error("Failed to unmarshal tenant", zap.Error(err))
				span.SetStatus(codes.Error, "Failed to unmarshal tenant")

				return err
			}

			span.SetAttributes(attribute.String("tenant_id", tent.ID.String()))
			log.Info("Tenant created", zap.String("tenant_id", tent.ID.String()), zap.String("tenant_name", tent.Name))
			span.SetStatus(codes.Ok, "Tenant created")
			return nil
		})

		if err != nil {
			log.Error("Failed to subscribe to tenant.created", zap.Error(err))
		}

	}()
}
