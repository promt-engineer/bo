package queue

import (
	"backoffice/internal/dto"
	"backoffice/pkg/rabbitmq"
	"backoffice/pkg/rabbitmq/consumer"
	"backoffice/pkg/rabbitmq/publisher"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
)

type Queue struct {
	cfg    *Config
	source *amqp.Connection

	router *Router

	publisherMap map[string]*publisher.Publisher
	listenerMap  map[string]*consumer.Consumer

	mu      *sync.RWMutex
	onClose chan *amqp.Error

	startWg   sync.WaitGroup
	startOnce sync.Once
}

var (
	ErrCanNotFindPublisher = errors.New("can not find publisher")

	contentType = "application/json"
)

func NewQueue(cfg *Config, handlers ...Handler) *Queue {
	q := &Queue{
		cfg: cfg,

		router: NewRouter(),

		publisherMap: map[string]*publisher.Publisher{},
		listenerMap:  map[string]*consumer.Consumer{},

		mu: &sync.RWMutex{},
	}

	for _, handler := range handlers {
		handler.Register(q.router)
	}

	return q
}

func (q *Queue) AddHandlers(handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(q.router)
	}
}

func (q *Queue) WaitTillStart() {
	q.startOnce.Do(func() {
		q.startWg.Add(1)

		go func() {
			for {
				err := q.connect(func() {
					zap.S().Info("Queue started")

					q.startWg.Done()
				})

				zap.S().Error(err)
			}
		}()

		q.startWg.Wait()
	})
}

func (q *Queue) connect(onStart func()) error {
	var err error

	q.source, err = rabbitmq.NewRabbitMQConn(q.cfg.Host, zap.L(), rabbitmq.NewOptions())
	if err != nil {
		return err
	}

	q.onClose = q.source.NotifyClose(make(chan *amqp.Error))

	for publisherName, publisherConf := range q.cfg.Publishers {
		q.publisherMap[publisherName], err = publisher.NewPublisher(q.source, zap.L(), publisher.NewOptions(
			publisher.WithExchangeName(publisherConf.ExchangeName()),
			publisher.WithBindingKey(publisherConf.BindingKey),
		))

		if err != nil {
			return err
		}
	}

	for listenerName, listenerConf := range q.cfg.Listeners {
		zap.S().Infof("Creating listener %s, binding key %s", listenerName, listenerConf.BindingKey)

		for i := 1; i <= listenerConf.Count; i++ {
			q.listenerMap[fmt.Sprintf("%s:%d", listenerName, i)], err = consumer.NewConsumer(q.source, zap.L(), consumer.NewOptions(
				consumer.WithExchangeName(listenerConf.ExchangeName()),
				consumer.WithBindingKey(listenerConf.BindingKey),
				consumer.WithQueueName(listenerConf.QueueName()),
				consumer.WithExchangeKind(listenerConf.ExchangeKind),
				consumer.WithConsumeNoWait(listenerConf.ExchangeNoWait),
				consumer.WithQueueAutoDelete(listenerConf.AutoDelete),
				consumer.WithQueueDurable(listenerConf.Durable),
			))
		}

		if err != nil {
			return err
		}
	}

	for _, cons := range q.listenerMap {
		go func(c *consumer.Consumer) {
			zap.S().Info("Start consuming")

			if err := c.StartConsumer(q.worker); err != nil {
				zap.S().Error(err)
			}
		}(cons)
	}

	onStart()

	return <-q.onClose
}

func (q *Queue) Send(ctx context.Context, publisherName string, msgType string, payload interface{}) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	pub, ok := q.publisherMap[publisherName]
	if !ok {
		return ErrCanNotFindPublisher
	}

	body, err := json.Marshal(dto.Msg{Type: msgType, Payload: payload})
	if err != nil {
		return err
	}

	zap.S().Info("start publish config to lord")
	err = pub.Publish(ctx, body, contentType)
	if errors.Is(err, amqp.ErrClosed) {
		zap.S().Error(err)
	}
	zap.S().Info("finished publish config to lord")

	return err
}

func (q *Queue) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for {
		select {
		case delivery, ok := <-messages:
			if !ok {
				return
			}

			var (
				err error
				msg dto.Msg
			)

			if err = json.Unmarshal(delivery.Body, &msg); err != nil {
				zap.S().Error(err)

				return
			}

			hf, ok := q.router.find(msg.Type)

			if !ok {
				zap.S().Errorf("can not find message type %v", msg.Type)
			}

			payload, err := json.Marshal(msg.Payload)
			if err != nil {
				zap.S().Error(err)

				return
			}

			hf(payload)

			if err = delivery.Ack(false); err != nil {
				zap.S().Error(err)

				return
			}

		case <-ctx.Done():
			return
		}
	}
}
