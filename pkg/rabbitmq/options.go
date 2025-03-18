package rabbitmq

import "time"

type Options struct {
	retryTimes int64
	backOff    time.Duration
}

type Option func(*Options)

func WithRetryTimes(retryTimes int64) Option {
	return func(o *Options) {
		o.retryTimes = retryTimes
	}
}

func WithBackOff(backOffSeconds time.Duration) Option {
	return func(o *Options) {
		o.backOff = 2 * time.Second
	}
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		retryTimes: retryTimes,
		backOff:    backOff,
	}

	for _, o := range opts {
		o(options)
	}

	return options
}
