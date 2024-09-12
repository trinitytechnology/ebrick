package tenant

import (
	"github.com/trinitytechnology/ebrick/module"
	"go.uber.org/zap"
)

type TenantModule struct {
}

// Install implements plugin.Plugin.
func (p *TenantModule) Initialize(opt *module.Options) error {
	// Init Tables
	log := opt.Logger
	log.Info("Initializing Tenant Module")

	// Migrate Tables
	opt.Database.AutoMigrate(&Tenant{})

	// Init Repository
	repo := NewRepository(opt.Database)

	// Init Service
	svc := NewTenantService(repo, opt.EventStream, opt.Logger)

	// Init Handler
	handler := NewTenantHandler(svc, opt.Logger)

	// setup routes
	setupApiRoutes(handler, opt.Router)

	// setup stream
	err := opt.EventStream.CreateStream("tenant", []string{"tenant.>"})

	if err != nil {
		log.Error("Failed to create stream", zap.Error(err))
		return err
	}
	log.Info("Tenant Module Initialized")
	return err
}

func (p *TenantModule) Name() string {
	return "tenant"
}

func (p *TenantModule) Version() string {
	return "1.0"
}

func (p *TenantModule) Description() string {
	return "Tenant Management"
}

func (p *TenantModule) Id() string {
	return "Tenant"
}

var Module TenantModule
