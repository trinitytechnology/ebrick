package cache

import (
	"time"

	"github.com/trinitytechnology/ebrick/config"
)

// Option represents a store option function.
type Option func(o *Options)

// Options represents the options for the cache store.
type Options struct {
	SynchronousSet            bool
	Cost                      int64
	Expiration                time.Duration
	Tags                      []string
	ClientSideCacheExpiration time.Duration

	Addrs    string
	User     string
	Password string
	Type     string
	Enable   bool
}

func newOptions(opts ...Option) *Options {
	cfg := config.GetConfig().Cache
	opt := &Options{
		Addrs:    cfg.Addrs,
		User:     cfg.User,
		Password: cfg.Password,
		Type:     cfg.Type,
		Enable:   cfg.Enable,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// config options functions
func Addrs(addrs string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func User(user string) Option {
	return func(o *Options) {
		o.User = user
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func Type(t string) Option {
	return func(o *Options) {
		o.Type = t
	}
}

func Enable(enable bool) Option {
	return func(o *Options) {
		o.Enable = enable
	}
}

// IsEmpty checks if the Options is empty.
func (o *Options) IsEmpty() bool {
	return o.Cost == 0 && o.Expiration == 0 && len(o.Tags) == 0
}

// WithExpiration allows specifying an expiration time when setting a value.
func WithExpiration(expiration time.Duration) Option {
	return func(o *Options) {
		o.Expiration = expiration
	}
}

// WithTags allows specifying associated tags to the current value.
func WithTags(tags []string) Option {
	return func(o *Options) {
		o.Tags = tags
	}
}

// WithClientSideCaching allows setting the client side caching, enabled by default.
// Currently to be used by Rueidis(redis) library only.
func WithClientSideCaching(clientSideCacheExpiration time.Duration) Option {
	return func(o *Options) {
		o.ClientSideCacheExpiration = clientSideCacheExpiration
	}
}

// InvalidateOption represents a cache invalidation function.
type InvalidateOption func(o *InvalidateOptions)

// InvalidateOptions represents the options for cache invalidation.
type InvalidateOptions struct {
	Tags []string
}

// isEmpty checks if the InvalidateOptions is empty.
func (o *InvalidateOptions) isEmpty() bool {
	return len(o.Tags) == 0
}

// ApplyInvalidateOptionsWithDefault applies the invalidate options with default values.
func ApplyInvalidateOptionsWithDefault(defaultOptions *InvalidateOptions, opts ...InvalidateOption) *InvalidateOptions {
	returnedOptions := ApplyInvalidateOptions(opts...)

	if returnedOptions.isEmpty() {
		returnedOptions = defaultOptions
	}

	return returnedOptions
}

// ApplyInvalidateOptions applies the invalidate options.
func ApplyInvalidateOptions(opts ...InvalidateOption) *InvalidateOptions {
	o := &InvalidateOptions{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithInvalidateTags allows setting the invalidate tags.
func WithInvalidateTags(tags []string) InvalidateOption {
	return func(o *InvalidateOptions) {
		o.Tags = tags
	}
}
