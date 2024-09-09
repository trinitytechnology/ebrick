package messaging

import (
	"context"
	"errors"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

var log *zap.Logger

func initNats(opt *Options) (*nats.Conn, nats.JetStreamContext) {
	log = logger.DefaultLogger
	opts := []nats.Option{
		nats.Name("NATS JetStream"),
		nats.UserInfo(opt.UserName, opt.Password),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2 * time.Second),
	}

	log.Info("Connecting to NATS", zap.String("url", opt.Url))
	nc, err := nats.Connect(opt.Url, opts...)
	if err != nil {
		log.Fatal("failed to connect to NATS", zap.Error(err))
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("failed to connect to JetStream", zap.Error(err))
	}

	log.Info("Connected to NATS JetStream")
	return nc, js
}

type natsJetStream struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	subs []*nats.Subscription
}

// CreateStream creates a JetStream stream with the specified name and subjects.
func (n *natsJetStream) CreateStream(stream string, subjects []string) error {
	_, err := n.js.AddStream(&nats.StreamConfig{
		Name:     stream,
		Subjects: subjects,
		Storage:  nats.FileStorage,
	})
	return err
}

// CreateConsumerGroup creates a JetStream consumer with the specified configuration.
func (n *natsJetStream) CreateConsumerGroup(stream, name string, config ConsumerConfig) error {
	_, err := n.js.AddConsumer(stream, &nats.ConsumerConfig{
		Durable:        name,
		AckWait:        config.AckWait,
		MaxDeliver:     config.MaxDeliver,
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverPolicy:  nats.DefaultPubRetryAttempts,
		BackOff:        config.BackOff,
		DeliverGroup:   config.DeliverGroup,
		DeliverSubject: config.DeliverSubject,
	})
	return err
}

// Close unsubscribes from all JetStream subscriptions and closes the NATS connection.
func (n *natsJetStream) Close() error {
	var errs []error
	for _, sub := range n.subs {
		if err := sub.Unsubscribe(); err != nil {
			errs = append(errs, err)
		}
	}
	n.conn.Close()

	if len(errs) > 0 {
		return errors.New("failed to close all subscriptions or connection")
	}
	return nil
}

// Publish publishes a CloudEvent to a JetStream subject with tracing context.
func (n *natsJetStream) Publish(ctx context.Context, ev event.Event) error {
	data, err := ev.MarshalJSON()
	if err != nil {
		log.Error("failed to marshal event", zap.Error(err))
		return err
	}
	headers := nats.Header{}

	// Check if tracing is enabled
	cfg := config.GetConfig().Observability
	if cfg.Tracing.Enable {
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	}

	_, err = n.js.PublishMsg(&nats.Msg{
		Subject: ev.Type(),
		Data:    data,
		Header:  headers,
	})
	return err
}

// Subscribe subscribes to a JetStream subject and processes incoming CloudEvents with the provided handler.
func (n *natsJetStream) Subscribe(subject, group string, handler func(ev *event.Event, ctx context.Context) error) error {
	// Check if the group parameter is empty
	if group == "" {
		return errors.New("group cannot be empty")
	} else {
		sub, err := n.js.QueueSubscribe(subject, group, func(msg *nats.Msg) {

			ctx := context.Background()

			// Check if tracing is enabled
			cfg := config.GetConfig().Observability
			if cfg.Tracing.Enable {
				ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Header))
			}

			var ev event.Event
			if err := ev.UnmarshalJSON(msg.Data); err != nil {
				log.Error("failed to unmarshal event", zap.Error(err))
				msg.Nak()
				return
			}

			if err := handler(&ev, ctx); err != nil {
				log.Error("failed to process event", zap.Error(err))
				msg.Nak()
				return
			}
			msg.Ack()
		}, nats.Durable(group), nats.ManualAck())

		if err != nil {
			log.Error("failed to subscribe to NATS JetStream", zap.Error(err))
			return err
		}
		log.Info("Successfully subscribed to subject", zap.String("subject", subject), zap.String("group", group))
		n.subs = append(n.subs, sub)
		return nil
	}

}

// SubscribeDLQ implements CloudEventStream.
func (n *natsJetStream) SubscribeDLQ(subject string, handler func(msg any, ctx context.Context) error) error {
	log.Info("Subscribing to NATS JetStream", zap.String("subject", subject))
	_, err := n.js.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(msg.Data, nil); err != nil {
			log.Error("failed to process event", zap.Error(err))
			msg.Nak()
			return
		}
		msg.Ack()
	}, nats.ManualAck())

	if err != nil {
		log.Error("failed to subscribe to NATS JetStream", zap.Error(err))
		return err
	}
	return nil
}

func NewNatsJetStream(opts ...Option) CloudEventStream {
	opt := newOptions(opts...)
	conn, js := initNats(opt)
	return &natsJetStream{
		conn: conn,
		js:   js,
	}
}
