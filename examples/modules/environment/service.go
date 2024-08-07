package environment

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type EnvironmentService interface {
	CreateEnvironment(ctx context.Context, tent Environment) (*Environment, error)
}
type envService struct {
	repo EnvironmentRepository
}

// CreateEnvironment implements EnvironmentService.
func (t *envService) CreateEnvironment(ctx context.Context, env Environment) (*Environment, error) {
	var tracer = otel.Tracer("env")
	_, span := tracer.Start(ctx, "Create Environment")
	span.SetAttributes(attribute.String("env.tenant_id", env.TenantId.String()), attribute.String("env.name", env.Name))
	defer span.End()

	createdEnvironment, err := t.repo.Create(env)
	if err != nil {
		log.Error("Error creating env", zap.Error(err))
		return nil, err
	}
	return createdEnvironment, nil
}

func NewEnvironmentService(repo EnvironmentRepository) EnvironmentService {
	return &envService{
		repo: repo,
	}
}
