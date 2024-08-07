package ebrick

import (
	"context"
	"time"

	"github.com/linkifysoft/ebrick/module"
	"go.uber.org/zap"
)

type App interface {
	Project() string
	Name() string
	Version() string
	Options() *Options
	RegisterModules(m ...module.Module) error
	Start() error
}

type application struct {
	opts *Options
	mm   *module.ModuleManager
}

// Version implements App.
func (a *application) Version() string {
	return a.opts.Version
}

// Project implements App.
func (a *application) Project() string {
	return a.opts.ProjectName
}

// Name implements App.
func (a *application) Name() string {
	return a.opts.ServiceName
}

// Options implements App.
func (a *application) Options() *Options {
	return a.opts
}

func NewApplication(opts ...Option) App {
	op := newOptions(opts...)

	mm := module.NewModuleManager(
		module.Logger(op.Logger),
		module.Database(op.Database),
		module.Cache(op.Cache),
		module.EventStream(op.EventStream),
		module.Router(op.HttpServer.GetRouter()),
	)

	return &application{
		opts: op,
		mm:   mm,
	}
}

// Start implements App.
func (a *application) Start() error {
	log := a.opts.Logger
	defer func() {
		// Increase timeout for tracer provider shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := a.opts.TracerProvider.Shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	err := a.opts.HttpServer.Start()

	return err
}

// RegisterModule registers a module.
func (a *application) RegisterModules(m ...module.Module) error {
	for _, module := range m {
		err := a.mm.RegisterModule(module)
		if err != nil {
			return err
		}
	}
	return nil
}
