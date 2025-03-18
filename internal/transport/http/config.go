package http

import "time"

type Config struct {
	Domain       string
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
