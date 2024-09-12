package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidiscompat"
	"github.com/trinitytechnology/ebrick/logger"
	"go.uber.org/zap"
)

const (
	RueidisType       = "rueidis"
	RueidisTagPattern = "gocache_tag_%s"

	defaultClientSideCacheExpiration = 10 * time.Second
)

// initRedisClient initializes the Redis client.
func initRedisClient(opts *Options) rueidis.Client {
	log := logger.DefaultLogger
	log.Info("Connecting to Redis", zap.String("url", opts.Addrs))
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{opts.Addrs}})
	if err != nil {
		log.Fatal("Failed to initialize Redis client", zap.Error(err))
	}
	log.Info("Connected to Redis")
	return client
}

// RedisStore is a Redis implementation of the Cache interface.
type RedisStore struct {
	client rueidis.Client
	opts   *Options
}

// NewRedisStore creates a new RedisStore instance.
func NewRedisStore(options ...Option) Cache {
	opts := newOptions(options...)
	client := initRedisClient(opts)

	if opts.ClientSideCacheExpiration == 0 {
		opts.ClientSideCacheExpiration = defaultClientSideCacheExpiration
	}

	return &RedisStore{
		client: client,
		opts:   opts,
	}
}

// HGet retrieves the value of a hash field.
func (store *RedisStore) HGet(ctx context.Context, key any, field any, options ...Option) (any, error) {
	cmd := store.client.B().Hget().Key(key.(string)).Field(field.(string)).Cache()
	res := store.client.DoCache(ctx, cmd, store.opts.ClientSideCacheExpiration)
	str, err := res.ToString()
	if rueidis.IsRedisNil(err) {
		err = NotFoundWithCause(err)
	}
	return str, err
}

// HGetAll retrieves all the fields and values from a hash.
func (store *RedisStore) HGetAll(ctx context.Context, key any, options ...Option) (map[string]any, error) {
	cmd := store.client.B().Hgetall().Key(key.(string)).Cache()
	res := store.client.DoCache(ctx, cmd, store.opts.ClientSideCacheExpiration)
	m, err := res.ToMap()
	if err != nil {
		return nil, err
	}
	data := make(map[string]any)
	for k, v := range m {
		if value, err := v.ToString(); err == nil {
			data[k] = value
		}
	}
	return data, nil
}

// HSet sets the value of a hash field.
func (store *RedisStore) HSet(ctx context.Context, key any, field any, value any, options ...Option) error {
	cmd := store.client.B().Hset().Key(key.(string)).FieldValue().FieldValue(field.(string), value.(string)).Build()
	err := store.client.Do(ctx, cmd).Error()
	if err != nil {
		return err
	}

	if tags := store.opts.Tags; len(tags) > 0 {
		store.setTags(ctx, key, tags)
	}
	ttl := int64(store.opts.Expiration)
	if ttl > 0 {
		expireResp := store.client.Do(ctx, store.client.B().Expire().Key(key.(string)).Seconds(ttl).Build())
		if expireResp.Error() != nil {
			return fmt.Errorf("failed to set expiration time for Hset: %w", expireResp.Error())
		}
	}
	return nil
}

// Get retrieves the value of a key.
func (store *RedisStore) Get(ctx context.Context, key any) (any, error) {
	cmd := store.client.B().Get().Key(key.(string)).Cache()
	res := store.client.DoCache(ctx, cmd, store.opts.ClientSideCacheExpiration)
	str, err := res.ToString()
	if rueidis.IsRedisNil(err) {
		err = NotFoundWithCause(err)
	}
	return str, err
}

// GetWithTTL retrieves the value of a key along with its time-to-live (TTL).
func (store *RedisStore) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	cmd := store.client.B().Get().Key(key.(string)).Cache()
	res := store.client.DoCache(ctx, cmd, store.opts.ClientSideCacheExpiration)
	str, err := res.ToString()
	if rueidis.IsRedisNil(err) {
		err = NotFoundWithCause(err)
	}
	return str, time.Duration(res.CacheTTL()) * time.Second, err
}

// Set sets the value of a key.
func (store *RedisStore) Set(ctx context.Context, key any, value any, options ...Option) error {
	opts := newOptions(options...)
	ttl := int64(opts.Expiration.Seconds())

	var cmd rueidis.Completed
	if ttl > 0 {
		cmd = store.client.B().Set().Key(key.(string)).Value(value.(string)).ExSeconds(ttl).Build()
	} else {
		cmd = store.client.B().Set().Key(key.(string)).Value(value.(string)).Build()
	}

	err := store.client.Do(ctx, cmd).Error()
	if err != nil {
		return err
	}

	if tags := opts.Tags; len(tags) > 0 {
		store.setTags(ctx, key, tags)
	}

	return nil
}

// setTags sets the tags for a key.
func (store *RedisStore) setTags(ctx context.Context, key any, tags []string) {
	ttl := 720 * time.Hour
	for _, tag := range tags {
		tagKey := fmt.Sprintf(RueidisTagPattern, tag)
		store.client.DoMulti(ctx,
			store.client.B().Sadd().Key(tagKey).Member(key.(string)).Build(),
			store.client.B().Expire().Key(tagKey).Seconds(int64(ttl.Seconds())).Build(),
		)
	}
}

// Delete deletes a key.
func (store *RedisStore) Delete(ctx context.Context, key any) error {
	return store.client.Do(ctx, store.client.B().Del().Key(key.(string)).Build()).Error()
}

// Invalidate invalidates the cache based on the provided options.
func (store *RedisStore) Invalidate(ctx context.Context, options ...InvalidateOption) error {
	opts := ApplyInvalidateOptions(options...)

	if tags := opts.Tags; len(tags) > 0 {
		for _, tag := range tags {
			tagKey := fmt.Sprintf(RueidisTagPattern, tag)

			cacheKeys, err := store.client.Do(ctx, store.client.B().Smembers().Key(tagKey).Build()).AsStrSlice()
			if err != nil {
				continue
			}

			for _, cacheKey := range cacheKeys {
				store.Delete(ctx, cacheKey)
			}

			store.Delete(ctx, tagKey)
		}
	}

	return nil
}

// GetType returns the type of the cache store.
func (store *RedisStore) GetType() string {
	return RueidisType
}

// Clear clears the cache.
func (store *RedisStore) Clear(ctx context.Context) error {
	return rueidiscompat.NewAdapter(store.client).FlushAll(ctx).Err()
}
