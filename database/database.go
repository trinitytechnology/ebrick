package database

import (
	"fmt"
	"sync"

	"github.com/trinitytechnology/ebrick/config"
	"github.com/trinitytechnology/ebrick/database/postgresql"
	"github.com/trinitytechnology/ebrick/logger"
	"gorm.io/gorm"
)

var (
	DefaultDataSource *gorm.DB = NewDatabase()
	once              sync.Once
)

// NewDatabase initializes the database connection based on the configuration.
func NewDatabase() *gorm.DB {
	logger := logger.DefaultLogger
	var db *gorm.DB
	once.Do(func() {
		cfg := config.GetConfig().Database
		if cfg.Type != "" && cfg.Enable {
			switch cfg.Type {
			case "postgresql":
				db = postgresql.InitDB()
			default:
				logger.Fatal(fmt.Sprintf("Database type %s is not supported", cfg.Type))
			}
		}
	})
	return db
}
