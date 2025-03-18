package consumer

import (
	"context"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type worker func(ctx context.Context, messages <-chan amqp.Delivery)

type Consumer struct {
	conn   *amqp.Connection
	logger *zap.Logger
	opts   *Options
}

func NewConsumer(conn *amqp.Connection, logger *zap.Logger, opts *Options) (*Consumer, error) {
	sub := &Consumer{
		conn:   conn,
		logger: logger,
		opts:   opts,
	}

	return sub, nil
}

func (c *Consumer) createChannel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	c.logger.Sugar().Infof("Declaring exchange: %s", c.opts.exchangeName)
	err = ch.ExchangeDeclare(
		c.opts.exchangeName,
		c.opts.exchangeKind,
		c.opts.exchangeDurable,
		c.opts.exchangeAutoDelete,
		c.opts.exchangeInternal,
		c.opts.exchangeNoWait,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		c.opts.queueName,
		c.opts.queueDurable,
		c.opts.queueAutoDelete,
		c.opts.queueExclusive,
		c.opts.queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	c.logger.Sugar().Infof("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		c.opts.exchangeName,
		c.opts.bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		c.opts.bindingKey,
		c.opts.exchangeName,
		c.opts.queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	c.logger.Sugar().Infof("Queue bound to exchange, starting to consume from queue, consumerTag: %v", c.opts.consumerTag)

	err = ch.Qos(
		c.opts.prefetchCount,
		c.opts.prefetchSize,
		c.opts.prefetchGlobal,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.Qos")
	}

	return ch, nil
}

func (c *Consumer) StartConsumer(fn worker) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.createChannel()
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		c.opts.queueName,
		c.opts.consumerTag,
		c.opts.consumeAutoAck,
		c.opts.consumeExclusive,
		c.opts.consumeNoLocal,
		c.opts.consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	forever := make(chan bool)

	for i := 0; i < c.opts.workerPoolSize; i++ {
		go fn(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Sugar().Error("ch.NotifyClose: %v", chanErr)
	<-forever

	return chanErr
}
