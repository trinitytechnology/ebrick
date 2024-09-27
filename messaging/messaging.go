package messaging

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/trinitytechnology/ebrick/logger"
	"go.uber.org/zap"
)

var (
	DefaultCloudEventStream CloudEventStream
	log                     = logger.DefaultLogger // zap.Logger
)

type CloudEventStream interface {
	Publish(topic string, ctx context.Context, ev event.Event) error
	Subscribe(topic, group string, handler func(msg *event.Event, ctx context.Context) error) error
	SubscribeDLQ(topic string, handler func(msg any, ctx context.Context) error) error
	CreateStream(stream string, topics []string) error
	CreateConsumerGroup(stream, name string, config ConsumerConfig) error
	Close() error
}

// init automatically initializes the CloudEventStream if the package is imported.
func init() {
	DefaultCloudEventStream = initializeCloudEventStream()
	if DefaultCloudEventStream == nil {
		log.Warn("Messaging is disabled or not properly configured")
	} else {
		log.Info("Messaging system initialized successfully")
	}
}

// initializeCloudEventStream sets up the appropriate CloudEventStream based on the config.
func initializeCloudEventStream() CloudEventStream {
	// Load the config
	err := loadConfig()
	if err != nil {
		return nil
	}

	if msgConfig == nil || !msgConfig.Enable {
		log.Error("Messaging configuration is not enabled or is missing")
		return nil
	}

	log.Info("Initializing messaging", zap.String("type", msgConfig.Type))

	switch msgConfig.Type {
	case "nats":
		return NewNatsJetStream()
	case "redis-stream":
		return NewRedisStream()
	default:
		log.Warn("Unsupported messaging type", zap.String("type", msgConfig.Type))
		return nil
	}
}
