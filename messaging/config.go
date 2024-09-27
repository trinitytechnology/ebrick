package messaging

import (
	"github.com/trinitytechnology/ebrick/config"
	"go.uber.org/zap"
)

// MessagingConfig represents the messaging configuration.
type MessagingConfig struct {
	Url      string
	UserName string
	Password string
	Enable   bool
	Type     string
}

var msgConfig *MessagingConfig // Accessible to the whole messaging package

// loadConfig loads the messaging configuration from the specified paths.
func loadConfig() error {
	cfg := struct {
		Messaging MessagingConfig
	}{}

	err := config.LoadConfig([]string{"."}, &cfg)
	if err != nil {
		log.Error("Error loading config", zap.Error(err))
		return err
	}

	msgConfig = &cfg.Messaging // Set the global configuration
	log.Info("Messaging configuration loaded successfully")
	return nil
}
