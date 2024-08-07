package ebrick

import (
	"github.com/linkifysoft/ebrick/cache"
	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/database"
	"github.com/linkifysoft/ebrick/logger"
	"github.com/linkifysoft/ebrick/messaging"
	"github.com/linkifysoft/ebrick/observability"
	"github.com/linkifysoft/ebrick/server"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Options struct {
	ProjectName string
	ServiceName string
	Version     string

	Database       *gorm.DB
	Cache          cache.Cache
	EventStream    messaging.CloudEventStream
	HttpServer     server.HttpServer
	TracerProvider *sdktrace.TracerProvider
	Logger         *zap.Logger
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	serviceCfg := config.GetConfig().Service

	opt := &Options{
		ProjectName:    serviceCfg.Project,
		ServiceName:    serviceCfg.Name,
		Version:        serviceCfg.Version,
		Database:       database.DefaultDataSource,
		Cache:          cache.DefaultCache,
		EventStream:    messaging.DefaultCloudEventStream,
		HttpServer:     server.DefaultServer,
		TracerProvider: observability.DefaultTraceProvider,
		Logger:         logger.DefaultLogger,
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func ProjectName(projectName string) Option {
	return func(o *Options) {
		o.ProjectName = projectName
	}
}

func ServiceName(serviceName string) Option {
	return func(o *Options) {
		o.ServiceName = serviceName
	}
}
