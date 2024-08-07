package messaging

import (
	"github.com/linkifysoft/ebrick/config"
)

type Options struct {
	Url      string
	UserName string
	Password string
	Enable   bool
	Type     string
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	cfg := config.GetConfig().Messaging
	opt := &Options{
		Url:      cfg.Url,
		UserName: cfg.UserName,
		Password: cfg.Password,
		Enable:   cfg.Enable,
		Type:     cfg.Type,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// option functions
func Url(url string) Option {
	return func(o *Options) {
		o.Url = url
	}
}

func UserName(userName string) Option {
	return func(o *Options) {
		o.UserName = userName
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func Enable(enable bool) Option {
	return func(o *Options) {
		o.Enable = enable
	}
}

func Type(t string) Option {
	return func(o *Options) {
		o.Type = t
	}
}
