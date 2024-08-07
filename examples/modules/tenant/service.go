package tenant

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/linkifysoft/ebrick/logger"
	"github.com/linkifysoft/ebrick/messaging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type TenantService interface {
	CreateTenant(ctx context.Context, tent Tenant) (*Tenant, error)
}
type tenantService struct {
	repo TenantRepository
	ces  messaging.CloudEventStream
}

// CreateTenant implements TenantService.
func (t *tenantService) CreateTenant(ctx context.Context, tent Tenant) (*Tenant, error) {
	log := logger.DefaultLogger

	createdTenant, err := t.repo.Create(tent)
	if err != nil {
		log.Error("Error creating tenant", zap.Error(err))
		return nil, err
	}

	t.PublishEvent(ctx, createdTenant)
	return createdTenant, nil
}

func (t *tenantService) PublishEvent(ctx context.Context, tent *Tenant) error {

	tracer := otel.Tracer("tenant")
	_, span := tracer.Start(ctx, "Send Tenant Created Event")
	span.SetAttributes(attribute.String("tenant.id", tent.ID.String()), attribute.String("module", "tenant"))

	defer span.End()
	// Create a new CloudEvent
	ce := event.New()
	ce.SetType("tenant.created")
	ce.SetSource("tenant-service")
	ce.SetData("application/json", tent)
	ce.SetExtension("tenant_id", tent.ID.String())

	// Publish a CloudEvent
	if err := t.ces.Publish(ctx, ce); err != nil {
		log.Error("Error publishing CloudEvent", zap.Error(err))
		return err
	}
	return nil
}

func NewTenantService(repo TenantRepository, eventStream messaging.CloudEventStream) TenantService {
	return &tenantService{
		repo: repo,
		ces:  eventStream,
	}
}
