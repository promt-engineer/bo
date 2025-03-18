package rpc

import "time"

type Config struct {
	Host              string
	Port              uint16
	MaxProcessingTime time.Duration
}
