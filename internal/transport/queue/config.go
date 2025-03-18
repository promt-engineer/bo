package queue

import (
	"backoffice/pkg/rabbitmq"
	"backoffice/utils"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Listeners  map[string]*ListenerConfig
	Publishers map[string]*PublisherConfig
	Host       *rabbitmq.Config
	Options    *Options
}

type ListenerConfig struct {
	ExchangeBase     string
	ExchangeKind     string
	QueueBase        string
	BindingKey       string
	HashedQueueName  bool
	Durable          bool
	AutoDelete       bool
	ExchangeInternal bool
	ExchangeNoWait   bool
	Count            int

	hashOnce sync.Once
	hashVal  string
}

func (lc *ListenerConfig) ExchangeName() string {
	return fmt.Sprintf("%v:%v", lc.ExchangeBase, lc.ExchangeKind)
}

func (lc *ListenerConfig) QueueName() string {
	lc.hashOnce.Do(func() {
		lc.hashVal = fmt.Sprintf("%x", utils.Rand64())
	})

	if lc.HashedQueueName {
		return fmt.Sprintf("%v:%v", lc.QueueBase, lc.hashVal)
	}

	return fmt.Sprintf("%v", lc.QueueBase)
}

type PublisherConfig struct {
	ExchangeBase string
	ExchangeKind string
	BindingKey   string
}

func (pc *PublisherConfig) ExchangeName() string {
	return fmt.Sprintf("%v:%v", pc.ExchangeBase, pc.ExchangeKind)
}

type Options struct {
	RetryTimes int
	BackOff    time.Time
}
