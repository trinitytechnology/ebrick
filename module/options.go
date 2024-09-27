package module

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Options struct {
	Database *gorm.DB
	Logger   *zap.Logger
	Router   *gin.Engine
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	opt := &Options{}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

// Option functions
func Database(db *gorm.DB) Option {
	return func(o *Options) {
		o.Database = db
	}
}

func Logger(l *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Router(r *gin.Engine) Option {
	return func(o *Options) {
		o.Router = r
	}
}
