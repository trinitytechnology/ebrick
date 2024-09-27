package server

import (
	"github.com/gin-gonic/gin"
	"github.com/trinitytechnology/ebrick/config"
	"github.com/trinitytechnology/ebrick/logger"
	"github.com/trinitytechnology/ebrick/observability"
	"github.com/trinitytechnology/ebrick/web"
	"github.com/trinitytechnology/ebrick/web/middleware"
	"go.uber.org/zap"
)

type Options struct {
	Port   int
	Env    string
	Logger *zap.Logger
	Router *gin.Engine
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	serverCfg := config.GetConfig().Server
	envCfg := config.GetConfig().Env
	obsCfg := config.GetConfig().Observability
	oidcCfg := config.GetConfig().Oidc

	webRouter := web.InitRouter()

	if oidcCfg.Enable {
		middleware.InitOIDC(&oidcCfg)
	}

	if obsCfg.Tracing.Enable {
		webRouter.Use(observability.TracingMiddleware(), observability.LoggingWithTraceIDMiddleware())
	}

	opt := Options{
		Port:   serverCfg.Port,
		Env:    envCfg,
		Logger: logger.DefaultLogger,
		Router: webRouter,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
