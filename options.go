package ebrick

import (
	"github.com/trinitytechnology/ebrick/cache"
	"github.com/trinitytechnology/ebrick/config"
	"github.com/trinitytechnology/ebrick/database"
	"github.com/trinitytechnology/ebrick/logger"
	"github.com/trinitytechnology/ebrick/observability"
	"github.com/trinitytechnology/ebrick/server"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Options struct {
	Name           string
	Version        string
	Database       *gorm.DB
	Cache          cache.Cache
	HttpServer     server.HttpServer
	TracerProvider *sdktrace.TracerProvider
	Logger         *zap.Logger
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	serviceCfg := config.GetConfig().Service

	opt := &Options{
		Name:           serviceCfg.Name,
		Version:        serviceCfg.Version,
		Database:       database.DefaultDataSource,
		Cache:          cache.DefaultCache,
		HttpServer:     server.DefaultServer,
		TracerProvider: observability.DefaultTraceProvider,
		Logger:         logger.DefaultLogger,
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

func GetVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func GetName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}
