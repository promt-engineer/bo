package rabbitmq

import (
	"backoffice/pkg/rabbitmq/consumer"
	"backoffice/pkg/rabbitmq/publisher"
	"context"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"testing"
	"time"
)

const (
	exchangeName       = "publisher_test"
	exchangeBindingKey = "test"
	queueName          = "test_test"
	consumerTag        = "pub-test-consumer"
	messageText        = "test message"
	messageTypeName    = "test.published"
)

func TestRabbitMQ(t *testing.T) {
	doneCh := make(chan struct{})
	errCh := make(chan error)

	type args struct {
		ctx         context.Context
		body        []byte
		contentType string
		fn          func(ctx context.Context, messages <-chan amqp.Delivery)
	}

	logger, _ := zap.NewProduction()
	conn, err := NewRabbitMQConn(&Config{
		Host:     "0.0.0.0",
		Port:     "5672",
		Password: "rabbitmq",
	},
		logger,
		NewOptions(),
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		conn    *amqp.Connection
		logger  *zap.Logger
		args    *args
		wantErr bool
	}{
		{
			name:   "Publish Test",
			conn:   conn,
			logger: logger,
			args: &args{
				ctx:         context.Background(),
				body:        []byte(messageText),
				contentType: "text/plain",
				fn: func(ctx context.Context, messages <-chan amqp.Delivery) {
					ticker := time.NewTicker(5 * time.Second)

					for {
						select {
						case msg := <-messages:
							if string(msg.Body) == messageText {
								if err := msg.Ack(false); err != nil {
									errCh <- err
									return
								}
								doneCh <- struct{}{}
								return
							} else {
								errCh <- errors.New("undefined message txt")
								return
							}
						case <-ticker.C:
							errCh <- errors.New("do not retrieve any messages")
							return
						}
					}
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		cons, err := consumer.NewConsumer(tt.conn, tt.logger, consumer.NewOptions(
			consumer.WithExchangeName(exchangeName),
			consumer.WithBindingKey(exchangeBindingKey),
			consumer.WithQueueName(queueName),
			consumer.WithConsumerTag(consumerTag),
		))
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			if err := cons.StartConsumer(tt.args.fn); err != nil {
				tt.logger.Sugar().Fatal(err)
			}
		}()

		t.Run(tt.name, func(t *testing.T) {
			p, err := publisher.NewPublisher(tt.conn, tt.logger, publisher.NewOptions(
				publisher.WithExchangeName(exchangeName),
				publisher.WithBindingKey(exchangeBindingKey),
				publisher.WithMessageTypeName(messageTypeName),
			))

			if err != nil {
				t.Fatal(err)
			}
			if err := p.Publish(tt.args.ctx, tt.args.body, tt.args.contentType); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		for {
			select {
			case e := <-errCh:
				t.Fatal(e)
			case <-doneCh:
				return
			}
		}
	}
}
