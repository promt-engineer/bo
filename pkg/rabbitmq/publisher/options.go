package publisher

type Options struct {
	exchangeName    string
	bindingKey      string
	messageTypeName string
	mandatory       bool
	immediate       bool
}

type Option func(o *Options)

func WithExchangeName(exchangeName string) Option {
	return func(o *Options) {
		o.exchangeName = exchangeName
	}
}

func WithBindingKey(bindingKey string) Option {
	return func(o *Options) {
		o.bindingKey = bindingKey
	}
}

func WithMessageTypeName(messageTypeName string) Option {
	return func(o *Options) {
		o.messageTypeName = messageTypeName
	}
}

func WithMandatory(mandatory bool) Option {
	return func(o *Options) {
		o.mandatory = mandatory
	}
}

func WithImmediate(immediate bool) Option {
	return func(o *Options) {
		o.immediate = immediate
	}
}

func NewOptions(opts ...Option) *Options {
	options := &Options{}

	for _, o := range opts {
		o(options)
	}

	return options
}
