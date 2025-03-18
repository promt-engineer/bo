package overlord

import (
	context "context"
	"crypto/tls"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetIntegratorConfig(ctx context.Context, integrator, game string) (*GetIntegratorConfigOut, error)
	SaveParams(ctx context.Context, in *SaveParamsIn) (out *SaveParamsOut, err error)
}

type client struct {
	api OverlordClient
}

type Config struct {
	Host     string
	Port     string
	IsSecure bool
}

type OpenBetResponse struct {
	TransactionID string
	Balance       int64
}

type CloseBetResponse struct {
	Balance int64
}

func newClient(host, port string, isSecure bool) (OverlordClient, error) {
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

	return NewOverlordClient(conn), nil
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

func (o *client) GetIntegratorConfig(ctx context.Context, integrator, game string) (*GetIntegratorConfigOut, error) {
	res, err := o.api.GetIntegratorConfig(ctx, &GetIntegratorConfigIn{Integrator: integrator, Game: game})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *client) SaveParams(ctx context.Context, in *SaveParamsIn) (out *SaveParamsOut, err error) {
	zap.S().Info("repo: SaveParams starting...")

	res, err := o.api.SaveParams(ctx, in)
	if err != nil {
		zap.S().Errorf("save params error: %s", err.Error())

		return nil, mapError(err)
	}

	return res, nil
}
