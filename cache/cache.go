package cache

import (
	"context"
	"time"

	"github.com/linkifysoft/ebrick/config"
	"github.com/linkifysoft/ebrick/logger"
	"go.uber.org/zap"
)

var DefaultCache Cache = NewCache()

// Cache is the interface that wraps the cache.
// Cache represents a generic cache interface.
type Cache interface {
	// Get retrieves the value associated with the given key from the cache.
	// If the key does not exist, it returns an error.
	Get(ctx context.Context, key any) (any, error)

	// GetWithTTL retrieves the value associated with the given key from the cache,
	// along with its time-to-live (TTL) duration.
	// If the key does not exist, it returns an error.
	GetWithTTL(ctx context.Context, key any) (any, time.Duration, error)

	// Set sets the value associated with the given key in the cache.
	// It also accepts optional options to customize the behavior of the cache operation.
	// If an error occurs, it returns the error.
	Set(ctx context.Context, key any, value any, options ...Option) error

	// HSet sets the value of a field in a hash stored at the given key in the cache.
	// It also accepts optional options to customize the behavior of the cache operation.
	// If an error occurs, it returns the error.
	HSet(ctx context.Context, key, field, value any, options ...Option) error

	// HGet retrieves the value of a field in a hash stored at the given key from the cache.
	// It also accepts optional options to customize the behavior of the cache operation.
	// If the key or field does not exist, it returns an error.
	HGet(ctx context.Context, key, field any, options ...Option) (any, error)

	// HGetAll retrieves all the fields and values from a hash stored at the given key in the cache.
	// It also accepts optional options to customize the behavior of the cache operation.
	// If the key does not exist, it returns an error.
	HGetAll(ctx context.Context, key any, options ...Option) (map[string]any, error)

	// Delete deletes the value associated with the given key from the cache.
	// If the key does not exist, it returns an error.
	Delete(ctx context.Context, key any) error

	// Invalidate invalidates the cache based on the provided options.
	// It also accepts optional options to customize the behavior of the cache invalidation.
	// If an error occurs, it returns the error.
	Invalidate(ctx context.Context, options ...InvalidateOption) error

	// Clear clears the entire cache.
	// If an error occurs, it returns the error.
	Clear(ctx context.Context) error

	// GetType returns the type of the cache.
	GetType() string
}

// InitCache is a function that initializes a cache.
func NewCache() Cache {
	var c Cache
	logger := logger.DefaultLogger

	if config.GetConfig().Cache.Enable {
		if config.GetConfig().Cache.Type == "redis" {
			c = NewRedisStore()
		} else {
			logger.Fatal("Invalid cache type", zap.String("type", config.GetConfig().Cache.Type))
		}
	}
	return c
}
