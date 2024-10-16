package messaging

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/trinitytechnology/ebrick/config"
	"github.com/trinitytechnology/ebrick/logger"
	"go.uber.org/zap"
)

var (
	DefaultCloudEventStream CloudEventStream = NewCloudEventStream()
)

var log *zap.Logger

type CloudEventStream interface {
	Publish(topic string, ctx context.Context, ev event.Event) error
	Subscribe(topic, group string, handler func(msg *event.Event, ctx context.Context) error) error
	SubscribeDLQ(topic string, handler func(msg any, ctx context.Context) error) error
	CreateStream(stream string, topics []string) error
	CreateConsumerGroup(stream, name string, config ConsumerConfig) error
	Close() error
}

func NewCloudEventStream() CloudEventStream {
	log = logger.DefaultLogger

	// check messaging is enabled then check type is Nats then init Nats
	cfg := config.GetConfig().Messaging
	if cfg.Enable {
		if cfg.Type == "nats" {
			return NewNatsJetStream()
		}
		if cfg.Type == "redis-stream" {
			return NewRedisStream()
		}
	}
	return nil
}
