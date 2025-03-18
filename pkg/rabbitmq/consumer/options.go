package consumer

type Options struct {
	exchangeName   string
	bindingKey     string
	queueName      string
	consumerTag    string
	workerPoolSize int

	exchangeKind       string
	exchangeDurable    bool
	exchangeAutoDelete bool
	exchangeInternal   bool
	exchangeNoWait     bool

	queueDurable    bool
	queueAutoDelete bool
	queueExclusive  bool
	queueNoWait     bool

	prefetchCount  int
	prefetchSize   int
	prefetchGlobal bool

	consumeAutoAck   bool
	consumeExclusive bool
	consumeNoLocal   bool
	consumeNoWait    bool
}

type Option func(*Options)

func WithExchangeName(exchangeName string) Option {
	return func(o *Options) {
		o.exchangeName = exchangeName
	}
}

func WithQueueName(queueName string) Option {
	return func(o *Options) {
		o.queueName = queueName
	}
}

func WithBindingKey(bindingKey string) Option {
	return func(o *Options) {
		o.bindingKey = bindingKey
	}
}

func WithConsumerTag(consumerTag string) Option {
	return func(o *Options) {
		o.consumerTag = consumerTag
	}
}

func WithWorkerPoolSize(workerPoolSize int) Option {
	return func(o *Options) {
		o.workerPoolSize = workerPoolSize
	}
}

func WithExchangeKind(exchangeKind string) Option {
	return func(o *Options) {
		o.exchangeKind = exchangeKind
	}
}

func WithExchangeDurable(exchangeDurable bool) Option {
	return func(o *Options) {
		o.exchangeDurable = exchangeDurable
	}
}

func WithExchangeAutoDelete(exchangeAutoDelete bool) Option {
	return func(o *Options) {
		o.exchangeAutoDelete = exchangeAutoDelete
	}
}

func WithExchangeInternal(exchangeInternal bool) Option {
	return func(o *Options) {
		o.exchangeInternal = exchangeInternal
	}
}

func WithExchangeNoWait(exchangeNoWait bool) Option {
	return func(o *Options) {
		o.exchangeNoWait = exchangeNoWait
	}
}

func WithQueueDurable(queueDurable bool) Option {
	return func(o *Options) {
		o.queueDurable = queueDurable
	}
}

func WithQueueAutoDelete(queueAutoDelete bool) Option {
	return func(o *Options) {
		o.queueAutoDelete = queueAutoDelete
	}
}

func WithQueueExclusive(queueExclusive bool) Option {
	return func(o *Options) {
		o.queueExclusive = queueExclusive
	}
}

func WithQueueNoWait(queueNoWait bool) Option {
	return func(o *Options) {
		o.queueNoWait = queueNoWait
	}
}

func WithPrefetchCount(prefetchCount int) Option {
	return func(o *Options) {
		o.prefetchCount = prefetchCount
	}
}

func WithPrefetchSize(prefetchSize int) Option {
	return func(o *Options) {
		o.prefetchSize = prefetchSize
	}
}

func WithPrefetchGlobal(prefetchGlobal bool) Option {
	return func(o *Options) {
		o.prefetchGlobal = prefetchGlobal
	}
}

func WithConsumeAutoAck(consumeAutoAck bool) Option {
	return func(o *Options) {
		o.consumeAutoAck = consumeAutoAck
	}
}

func WithConsumeExclusive(consumeExclusive bool) Option {
	return func(o *Options) {
		o.consumeExclusive = consumeExclusive
	}
}

func WithConsumeNoLocal(consumeNoLocal bool) Option {
	return func(o *Options) {
		o.consumeNoLocal = consumeNoLocal
	}
}

func WithConsumeNoWait(consumeNoWait bool) Option {
	return func(o *Options) {
		o.consumeNoWait = consumeNoWait
	}
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		workerPoolSize:  10,
		exchangeKind:    "direct",
		exchangeDurable: true,
		queueDurable:    true,
		prefetchCount:   1,
		prefetchSize:    0,
	}

	for _, o := range opts {
		o(options)
	}

	return options
}
