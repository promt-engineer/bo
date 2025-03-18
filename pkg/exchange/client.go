package exchange

import (
	"context"
	"crypto/tls"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RateKey struct {
	From, To string
}

type RatesBag map[RateKey][]*RateItem

type Client interface {
	GetRates(ctx context.Context, from string, to []string, start, end time.Time) (RatesBag, error)
	UpdateCurrencies(ctx context.Context, currency, baseCurrency string) (*Status, error)
	GetCurrencyRates(ctx context.Context, in *AllCurrencyRatesIn) (*AllCurrencyRatesOut, error)
	AddCurrencyRate(ctx context.Context, from, to string, rate float64) (*CurrencyRates, error)
	DeleteCurrencyRate(ctx context.Context, from, to string, rate float64, createdAt time.Time) (*Status, error)
}

type Config struct {
	Host     string
	Port     string
	IsSecure bool
}

func NewClient(cfg *Config) (Client, error) {
	var err error

	service := &client{}
	service.api, err = newClient(cfg.Host, cfg.Port, cfg.IsSecure)

	if err != nil {
		return service, err
	}

	return service, nil
}

func newClient(host, port string, isSecure bool) (ExchangeServiceClient, error) {
	addr := host + ":" + port

	var (
		conn *grpc.ClientConn
		err  error
	)

	if isSecure {
		config := &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}

		conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err != nil {
		zap.S().Errorf("can not dial %v: %v", addr, err)

		return nil, err
	}

	return NewExchangeServiceClient(conn), nil
}

type client struct {
	api ExchangeServiceClient
}

func (c *client) GetRates(ctx context.Context, from string, to []string, start, end time.Time) (RatesBag, error) {
	out, err := c.api.GetRates(ctx, &RatesIn{
		From:  from,
		To:    to,
		Start: timestamppb.New(start),
		End:   timestamppb.New(end),
	})

	if err != nil {
		return nil, err
	}

	bag := RatesBag{}

	for key, value := range out.Bag {
		sp := strings.Split(key, "->")
		bag[RateKey{From: sp[0], To: sp[1]}] = value.Items
	}

	return bag, nil
}

func (c *client) UpdateCurrencies(ctx context.Context, currency, baseCurrency string) (*Status, error) {
	out, err := c.api.UpdateCurrencies(ctx, &UpdateCurrency{
		Currency:     currency,
		BaseCurrency: baseCurrency,
	})

	if err != nil {
		return &Status{}, err
	}

	return out, nil
}

func (c *client) GetCurrencyRates(ctx context.Context, in *AllCurrencyRatesIn) (*AllCurrencyRatesOut, error) {
	return c.api.GetAllCurrencyRates(ctx, in)
}

func (c *client) AddCurrencyRate(ctx context.Context, from, to string, rate float64) (*CurrencyRates, error) {
	out, err := c.api.AddCurrencyRate(ctx, &AddCurrencyRateIn{
		From: from,
		To:   to,
		Rate: rate,
	})
	log.Printf("client rpc response: %v", out)
	if err != nil {
		return &CurrencyRates{}, err
	}
	return out, nil
}

func (c *client) DeleteCurrencyRate(ctx context.Context, from, to string, rate float64, createdAt time.Time) (*Status, error) {
	out, err := c.api.DeleteCurrencyRate(ctx, &DeleteCurrencyRateIn{
		From:      from,
		To:        to,
		Rate:      rate,
		CreatedAt: timestamppb.New(createdAt),
	})
	if err != nil {
		return &Status{}, err
	}
	return out, nil
}
