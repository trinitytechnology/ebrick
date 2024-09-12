package environment

import (
	"context"

	"github.com/trinitytechnology/ebrick/observability"
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

	_, span := observability.StartSpan(ctx, Module.Name(), "Create Environment")
	span.SetAttributes(attribute.String("env_tenant_id", env.TenantId.String()), attribute.String("env_name", env.Name))
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
