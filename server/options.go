package server

import (
	"github.com/gin-gonic/gin"
	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"github.com/linkifysoft/ebrick/server/middleware"
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
	env := config.GetConfig().Env

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.TracingMiddleware(), middleware.LoggingMiddleware())
	setupProbeRoute(router)

	opt := Options{
		Port:   serverCfg.Port,
		Env:    env,
		Logger: logger.DefaultLogger,
		Router: router,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
