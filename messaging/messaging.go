package messaging

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/linkifysoft/ebrick/config"
)

var (
	DefaultCloudEventStream CloudEventStream = NewCloudEventStream()
)

type CloudEventStream interface {
	Publish(ctx context.Context, ev event.Event) error
	Subscribe(subject, group string, handler func(msg *event.Event, ctx context.Context) error) error
	SubscribeDLQ(subject string, handler func(msg any, ctx context.Context) error) error
	CreateStream(stream string, subjects []string) error
	CreateConsumerGroup(stream, name string, config ConsumerConfig) error
	Close() error
}

func NewCloudEventStream() CloudEventStream {
	// check messaging is enabled then check type is Nats then init Nats
	cfg := config.GetConfig().Messaging
	if cfg.Enable {
		if cfg.Type == "nats" {
			return NewNatsJetStream()
		}
	}
	return nil
}
