package history

import (
	"context"
	"crypto/tls"
	"math"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetAllSpins(ctx context.Context, in *FinancialBase) (*GetAllSpinsOut, error)
	AllGamingSession(ctx context.Context, in *FinancialBase) (*GetAllGameSessionsOut, error)
	GetSpin(ctx context.Context, roundID, currency string) (*GetSpinOut, error)
	GetSession(ctx context.Context, gameIDs []string, sessionID, currency string) (*GameSessionOut, error)

	GetAggregatedReportByGame(ctx context.Context, in *GetAggregatedReportFilters) ([]*GetAggregatedReportByGameItem, error)
	GetAggregatedReportByCountry(ctx context.Context, in *GetAggregatedReportFilters) ([]*GetAggregatedReportByCountryItem, error)
	GetFinancialReport(ctx context.Context, in *FinancialBase) (*FinancialReport, error)

	GetSpins(ctx context.Context, in *GetFinancialIn) (*GetSpinsOut, error)
	GetSessions(ctx context.Context, in *GetFinancialIn) (*GetSessionsOut, error)

	GetHosts(ctx context.Context, in *FinancialBase) ([]string, error)
	GetCurrencies(ctx context.Context, in *FinancialBase) ([]string, error)
	IntegratorOperatorsMap(ctx context.Context, games []string) (map[string][]string, error)
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

func newClient(host, port string, isSecure bool) (HistoryServiceClient, error) {
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

		conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	} else {
		conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	}

	if err != nil {
		zap.S().Errorf("can not dial %v: %v", addr, err)

		return nil, err
	}

	return NewHistoryServiceClient(conn), nil
}

type client struct {
	api HistoryServiceClient
}

func (c *client) GetAllSpins(ctx context.Context, in *FinancialBase) (*GetAllSpinsOut, error) {
	return c.api.GetAllSpins(ctx, in)
}

func (c *client) AllGamingSession(ctx context.Context, in *FinancialBase) (*GetAllGameSessionsOut, error) {
	return c.api.GetAllGameSession(ctx, in)
}

func (c *client) GetSpin(ctx context.Context, roundID, currency string) (*GetSpinOut, error) {
	return c.api.GetSpin(ctx, &GetSpinIn{RoundId: roundID, ConvertCurrency: currency})
}

func (c *client) GetSession(ctx context.Context, gameIDs []string, sessionID, currency string) (*GameSessionOut, error) {
	return c.api.GetSession(ctx, &GetSessionIn{SessionId: sessionID, Games: gameIDs, ConvertCurrency: currency})
}

func (c *client) GetHosts(ctx context.Context, in *FinancialBase) ([]string, error) {
	out, err := c.api.GetHosts(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Items, nil
}

func (c *client) GetCurrencies(ctx context.Context, in *FinancialBase) ([]string, error) {
	out, err := c.api.GetCurrencies(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Items, nil
}

func (c *client) IntegratorOperatorsMap(ctx context.Context, games []string) (map[string][]string, error) {
	out, err := c.api.GetIntegratorOperators(ctx, &GamesIn{Games: games})
	if err != nil {
		return nil, err
	}

	return lo.MapEntries(out.Map, func(key string, value *DictionaryOut) (string, []string) {
		return key, value.Items
	}), nil
}

func (c *client) GetAggregatedReportByGame(ctx context.Context, in *GetAggregatedReportFilters) ([]*GetAggregatedReportByGameItem, error) {
	out, err := c.api.GetAggregatedReportByGame(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Items, nil
}

func (c *client) GetAggregatedReportByCountry(ctx context.Context, in *GetAggregatedReportFilters) ([]*GetAggregatedReportByCountryItem, error) {
	out, err := c.api.GetAggregatedReportByCountry(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Items, nil
}

func (c *client) GetSpins(ctx context.Context, in *GetFinancialIn) (*GetSpinsOut, error) {
	return c.api.GetSpins(ctx, in)
}

func (c *client) GetSessions(ctx context.Context, in *GetFinancialIn) (*GetSessionsOut, error) {
	return c.api.GetSessions(ctx, in)
}

func (c *client) GetFinancialReport(ctx context.Context, in *FinancialBase) (*FinancialReport, error) {
	out, err := c.api.GetFinancialReport(ctx, in)
	if err != nil {
		return nil, err
	}

	return out.Report, nil
}
