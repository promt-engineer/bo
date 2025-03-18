package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

const (
	retryTimes = 5
	backOff    = 2 * time.Second

	dsn = "amqp://%s:%s@%s:%s/"
)

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

func NewRabbitMQConn(cfg *Config, logger *zap.Logger, opts *Options) (*amqp.Connection, error) {
	var (
		conn   *amqp.Connection
		counts int64
	)

	for {
		url := fmt.Sprintf(dsn, cfg.Username, cfg.Password, cfg.Host, cfg.Port)
		connection, err := amqp.Dial(url)
		if err != nil {
			logger.Sugar().Errorf("RabbitMq at %s not ready...\n", url)
			counts++
		} else {
			conn = connection

			break
		}

		if counts > opts.retryTimes {
			logger.Sugar().Error(err)

			return nil, ErrCannotConnectRabbitMQ
		}

		logger.Sugar().Infof("Backing off for %v...", opts.backOff)
		time.Sleep(opts.backOff)

		continue
	}

	logger.Info("Connected to RabbitMQ!")

	return conn, nil
}
