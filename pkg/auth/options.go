package auth

import "time"

type TokenOptions struct {
	ID           string
	Secret       string
	RefreshToken string
	Expiry       time.Duration
}

type TokenOption func(o *TokenOptions)

func WithExpiry(ex time.Duration) TokenOption {
	return func(o *TokenOptions) {
		o.Expiry = ex
	}
}

func WithTokenID(id string) TokenOption {
	return func(o *TokenOptions) {
		o.ID = id
	}
}

func WithCredentials(id, secret string) TokenOption {
	return func(o *TokenOptions) {
		o.ID = id
		o.Secret = secret
	}
}

func WithRefreshToken(rt string) TokenOption {
	return func(o *TokenOptions) {
		o.RefreshToken = rt
	}
}

func NewTokenOptions(opts ...TokenOption) TokenOptions {
	var options TokenOptions
	for _, o := range opts {
		o(&options)
	}

	if options.Expiry == 0 {
		options.Expiry = time.Minute
	}

	return options
}

type GenerateOptions struct {
	ID      string
	Subject string
}

type GenerateOption func(o *GenerateOptions)

func WithID(id string) GenerateOption {
	return func(o *GenerateOptions) {
		o.ID = id
	}
}

func WithSubject(s string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Subject = s
	}
}

func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	var options GenerateOptions
	for _, o := range opts {
		o(&options)
	}
	return options
}
