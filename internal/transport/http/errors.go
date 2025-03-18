package http

import "errors"

var (
	ErrAuthHeaderIsRequired = errors.New("auth header is required")
)
