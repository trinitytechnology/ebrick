package messaging

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/redis/rueidis"
	"github.com/trinitytechnology/ebrick/config"
	"github.com/trinitytechnology/ebrick/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

type redisStream struct {
	client           rueidis.Client
	ctx              context.Context
	consumer_configs map[string]ConsumerConfig
}

// DefaultConsumerConfig provides default values for ConsumerConfig.
var DefaultConsumerConfig = ConsumerConfig{
	MaxDeliver:     3,               // Default max delivery attempts
	AckWait:        2 * time.Second, // Default wait time for acknowledgment
	DeliverSubject: "",              // Default subject for delivery (can be configured as needed)
	// Add other default values as necessary...
}

// InitRedisClient initializes and returns a new Redis client.
func InitRedisClient() *rueidis.Client {
	cfg := config.GetConfig().Messaging

	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{cfg.Url}})
	if err != nil {
		log.Fatal("failed to create Redis client", zap.Error(err))
		return nil
	}
	return &client
}

// NewRedisStream creates a new instance of redisStream, connecting to Redis.
func NewRedisStream() CloudEventStream {
	redisURL := config.GetConfig().Messaging.Url
	client := InitRedisClient()
	log.Info("Connected to Redis", zap.String("url", redisURL))

	return &redisStream{
		client:           *client,
		ctx:              context.Background(),
		consumer_configs: make(map[string]ConsumerConfig),
	}
}

// CreateStream logs a warning that Redis does not require explicit stream creation.
func (r *redisStream) CreateStream(stream string, subjects []string) error {
	log.Warn("Redis no need to create stream first. Because we need to create default consumer group or publish a message to create a Stream", zap.String("stream", stream), zap.Strings("subjects", subjects))
	return nil
}

// CreateConsumerGroup creates a consumer group for the specified stream; returns error if it fails.
func (r *redisStream) CreateConsumerGroup(stream, group string, config ConsumerConfig) error {
	// Store the consumer config in the map
	if !utils.IsBlank(&config.GroupName) {
		r.consumer_configs[group] = config
	}

	startId := utils.Default(&config.StartID, "0")
	builder := r.client.B().XgroupCreate().Key(stream).Group(group).Id(startId).Mkstream()
	resp := r.client.Do(r.ctx, builder.Build())
	if resp.Error() != nil {
		if strings.Contains(resp.Error().Error(), "BUSYGROUP Consumer Group name already exists") {
			return nil // Consumer group already exists; nothing to do here
		}
		return fmt.Errorf("failed to create consumer group: %w", resp.Error())
	}

	return nil
}

// Close performs any necessary cleanup for the redisStream.
func (r *redisStream) Close() error {
	return nil
}

// Publish sends an event to the specified stream and adds tracing information if enabled.
func (r *redisStream) Publish(stream string, ctx context.Context, ev event.Event) error {
	data, err := ev.MarshalJSON()
	if err != nil {
		log.Error("failed to marshal event", zap.Error(err))
		return err
	}

	headers := make(map[string]string)
	cfg := config.GetConfig().Observability
	if cfg.Tracing.Enable {
		otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(headers))
	}

	builder := r.client.B().Xadd().Key(stream).Id("*").FieldValue().FieldValue("event", rueidis.BinaryString(data))

	if len(headers) > 0 {
		headersJSON := utils.MarshalJSON(headers)
		builder.FieldValue("trace", headersJSON)
	}

	resp := r.client.Do(r.ctx, builder.Build())
	if resp.Error() != nil {
		return fmt.Errorf("failed to add message to stream: %w", resp.Error())
	}
	return nil
}

// Subscribe sets up a consumer to process messages from a stream using the specified group.
func (r *redisStream) Subscribe(stream, group string, handler func(ev *event.Event, ctx context.Context) error) error {
	if group == "" {
		return errors.New("group cannot be empty")
	}

	// Fetch consumer configuration, use defaults if not found
	config, exists := r.consumer_configs[group]
	if !exists {
		log.Debug("Consumer config not found, using default values", zap.String("group", group))
		config = DefaultConsumerConfig
		config.GroupName = group
	}

	if err := r.CreateConsumerGroup(stream, group, config); err != nil {
		return fmt.Errorf("error creating consumer group: %w", err)
	}

	go func() {
		for {
			msgId, ev, err := r.ConsumeMessages(group, GenerateConsumerName(group), ">", 10, 0, stream)
			if err != nil {
				log.Error("Error consuming messages from stream", zap.Error(err))
				continue
			}

			attempts := 0
			for {
				err := handler(&ev, r.ctx)
				if err == nil {
					r.ackMsg(stream, group, msgId)
					break // Successful processing
				}

				attempts++
				log.Warn("Processing failed, attempting retry", zap.String("msgId", msgId), zap.Int("attempt", attempts))

				if attempts >= config.MaxDeliver {
					log.Error("Max retries exceeded, sending to DLQ", zap.String("msgId", msgId))
					r.PublishDLQ(ev, config.DeliverSubject)
					r.ackMsg(stream, group, msgId)
					break
				}
				time.Sleep(config.AckWait) // Use the configured wait time
			}
		}
	}()

	log.Info("Successfully subscribed to stream", zap.String("stream", stream), zap.String("group", group))
	return nil
}

// ConsumeMessages reads messages from a specified group and streams; returns message ID and event.
func (r *redisStream) ConsumeMessages(groupName, consumerName, startID string, count int64, block int64, streams ...string) (string, event.Event, error) {
	if len(streams) == 0 {
		return "", event.Event{}, fmt.Errorf("no streams specified")
	}

	streamIDs := make([]string, 0, len(streams)*2)
	for range streams {
		streamIDs = append(streamIDs, startID)
	}

	builder := r.client.B().Xreadgroup().Group(groupName, GenerateConsumerName(consumerName)).Block(block).Streams().Key(streams...).Id(streamIDs...)
	resp := r.client.Do(r.ctx, builder.Build())
	if resp.Error() != nil {
		return "", event.Event{}, fmt.Errorf("failed to read messages from stream: %w", resp.Error())
	}

	xEntry, err := resp.AsXRead()
	if err != nil {
		return "", event.Event{}, fmt.Errorf("failed to read messages from stream: response is not a valid XRead")
	}

	for _, entries := range xEntry {
		for _, fields := range entries {
			if data, ok := fields.FieldValues["event"]; ok {
				ev, err := utils.UnmarshalJSON[event.Event](data)
				if err != nil {
					log.Error("failed to unmarshal event", zap.Error(err))
					continue
				}

				if config.GetConfig().Observability.Tracing.Enable {
					if traceData, ok := fields.FieldValues["trace"]; ok {
						carrier, err := utils.UnmarshalJSON[map[string]string](traceData)
						if err == nil {
							r.ctx = otel.GetTextMapPropagator().Extract(r.ctx, propagation.MapCarrier(carrier))
						} else {
							log.Error("failed to unmarshal trace data", zap.Error(err))
						}
					}
				}

				return fields.ID, ev, nil // Return Redis message ID and event
			}
		}
	}

	return "", event.Event{}, fmt.Errorf("no data found in stream messages")
}

// SubscribeDLQ subscribes to a dead letter queue (DLQ) stream for processing.
func (r *redisStream) SubscribeDLQ(stream string, handler func(msg any, ctx context.Context) error) error {
	log.Info("Subscribing to Redis DLQ", zap.String("subject", stream))

	dlqGroup := stream + "-dlq-group"
	if err := r.CreateConsumerGroup(stream, dlqGroup, ConsumerConfig{StartID: "0"}); err != nil {
		log.Error("Error creating DLQ consumer group, it may already exist", zap.Error(err))
	}

	go func() {
		for {
			msgId, ev, err := r.ConsumeMessages(dlqGroup, dlqGroup, ">", 1, 0, stream)
			if err != nil {
				log.Error("Error consuming messages from DLQ stream", zap.Error(err))
				time.Sleep(time.Second) // Wait before retrying
				continue
			}

			if err := handler(ev.Data(), r.ctx); err != nil {
				log.Error("Failed to process DLQ event", zap.Error(err))
				continue
			}

			r.ackMsg(stream, dlqGroup, msgId) // Acknowledge the message
		}
	}()

	return nil
}

// ackMsg acknowledges the processing of a message in the specified stream and group.
func (r *redisStream) ackMsg(stream, group, messageID string) {
	builder := r.client.B().Xack().Key(stream).Group(group).Id(messageID)
	resp := r.client.Do(r.ctx, builder.Build())
	if resp.Error() != nil {
		log.Error("failed to acknowledge message", zap.Error(resp.Error()))
	}
}

// PublishDLQ sends an event to the dead letter queue (DLQ) stream.
func (r *redisStream) PublishDLQ(ev event.Event, dlqStream string) {
	data, _ := ev.MarshalJSON()
	r.client.B().Xadd().Key(dlqStream).Id("*").FieldValue().FieldValue("event", rueidis.BinaryString(data)).Build()
}
