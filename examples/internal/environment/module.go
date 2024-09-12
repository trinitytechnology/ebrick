package environment

import (
	"github.com/trinitytechnology/ebrick/module"
)

type EnvironmentModule struct {
}

// Install implements plugin.Plugin.
func (p *EnvironmentModule) Initialize(opt *module.Options) error {
	// Init Tables
	log := opt.Logger

	log.Info("Initializing Environment Module")

	// Migrate Tables
	opt.Database.AutoMigrate(&Environment{})

	// Init Repository
	repo := NewRepository(opt.Database)

	// Init Service
	svc := NewEnvironmentService(repo, opt.Logger)

	// Init Handler
	handler := NewEnvironmentHandler(svc, opt.Logger)

	// setup routes
	setupRoutes(handler, opt.Router)

	// setup stream
	configureConsumers(opt.EventStream, svc, opt.Logger)

	log.Info("Environment Module Initialized")
	return nil
}

func (p *EnvironmentModule) Name() string {
	return "env"
}

func (p *EnvironmentModule) Version() string {
	return "1.0"
}

func (p *EnvironmentModule) Description() string {
	return "Environment Management"
}

func (p *EnvironmentModule) Id() string {
	return "Environment"
}

var Module EnvironmentModule
