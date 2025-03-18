package totp

import "errors"

var (
	ErrValidateSecretInvalidBase32 = errors.New("Decoding of secret as base32 failed.")
	ErrValidateInputInvalidLength  = errors.New("Input length unexpected")
	ErrGenerateMissingUsername     = errors.New("Username must be set")
)

type GenerateOptions struct {
	username string
}

type GenerateOption func(options *GenerateOptions)

func WithUsername(username string) GenerateOption {
	return func(options *GenerateOptions) {
		options.username = username
	}
}

func NewGenerateOptions(opts ...GenerateOption) *GenerateOptions {
	options := &GenerateOptions{}

	for _, o := range opts {
		o(options)
	}

	return options
}
