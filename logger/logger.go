package logger

import (
	"log"

	"github.com/trinitytechnology/ebrick/config"
	"go.uber.org/zap"
)

var DefaultLogger *zap.Logger

func init() {
	environment := config.GetConfig().Env
	DefaultLogger = NewLogger(environment)
}

func NewLogger(env string) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	if env == "production" {
		config = zap.NewProductionConfig()
	}
	zapLogger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()
	zapLogger.Info("Starting application", zap.String("env", env))
	return zapLogger
}
