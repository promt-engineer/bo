package pgsql

import (
	"time"
)

type Config struct {
	Host              string
	Port              uint16
	Name              string
	User              string
	Pass              string
	Compression       string
	ConnectionTimeout time.Duration
	PingInterval      time.Duration
	MinConnections    int
	MaxConnections    int
}
