package postgresql

import (
	"fmt"

	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the PostgreSQL database connection and returns a *gorm.DB instance.
func InitDB() *gorm.DB {
	log := logger.DefaultLogger
	// Get the database configuration from the config package
	cfg := config.GetConfig()
	log.Info("Connecting to PostgreSQL database", zap.String("host", cfg.Database.Host), zap.String("dbname", cfg.Database.DBName), zap.String("sslmode", cfg.Database.SSLMode))
	// Create the connection string using the configuration values
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.DBName, cfg.Database.SSLMode, cfg.Database.Password)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	log.Info("Connected to PostgreSQL database")
	return db
}
