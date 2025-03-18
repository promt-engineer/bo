package publisher

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type Publisher struct {
	conn   *amqp.Connection
	logger *zap.Logger
	opts   *Options
}

func NewPublisher(conn *amqp.Connection, logger *zap.Logger, opts *Options) (*Publisher, error) {
	pub := &Publisher{
		conn:   conn,
		logger: logger,
		opts:   opts,
	}

	return pub, nil
}

func (p *Publisher) Publish(ctx context.Context, body []byte, contentType string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	p.logger.Sugar().Infof("Publishing message Exchange: %s, RoutingKey: %s", p.opts.exchangeName, p.opts.bindingKey)

	if err = ch.PublishWithContext(
		ctx,
		p.opts.exchangeName,
		p.opts.bindingKey,
		p.opts.mandatory,
		p.opts.immediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
			Type:         p.opts.messageTypeName,
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	return nil
}
