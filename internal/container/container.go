package container

import (
	"context"
	"fmt"
	"sync"

	"backoffice/internal/config"
	"backoffice/internal/constants"
	"backoffice/internal/services"
	"backoffice/internal/transport/http"
	"backoffice/internal/transport/rpc"
	"backoffice/pkg/auth/jwt"
	"backoffice/pkg/exchange"
	"backoffice/pkg/history"
	"backoffice/pkg/mailgun"
	"backoffice/pkg/overlord"
	"backoffice/pkg/pgsql"
	"backoffice/pkg/redis"
	"backoffice/pkg/totp"
	"backoffice/pkg/validator"
	"bitbucket.org/play-workspace/gocommon/tracer"

	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

var container di.Container
var once sync.Once

func Build(ctx context.Context, wg *sync.WaitGroup) di.Container {
	once.Do(func() {
		builder, _ := di.NewBuilder()
		defs := []di.Def{
			{
				Name: constants.LoggerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)
					c := zap.NewDevelopmentConfig()
					level, err := zap.ParseAtomicLevel(cfg.LogLevel)
					if err != nil {
						return nil, err
					}

					c.Level = level
					logger, err := c.Build()

					if err != nil {
						return nil, fmt.Errorf("can't initialize zap logger: %v", err)
					}

					zap.ReplaceGlobals(logger)

					return logger, nil
				},
			},
			{
				Name: constants.ConfigName,
				Build: func(ctn di.Container) (interface{}, error) {
					return config.New()
				},
			},
			{
				Name: constants.PgSQLConnectionName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return pgsql.NewPgSQLConnection(cfg.PgSQLConfig)
				},
			},
			{
				Name: constants.RedisName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return redis.New(cfg.RedisConfig)
				},
			},
			{
				Name: constants.ValidatorName,
				Build: func(ctn di.Container) (interface{}, error) {
					return validator.New()
				},
			},
			{
				Name: constants.HistoryName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return history.NewClient(cfg.HistoryConfig)
				},
			},
			{
				Name: constants.HTTPServerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)
					handlers := []http.Handler{
						ctn.Get(constants.MetaHTTPHandlerName).(http.Handler),
						ctn.Get(constants.PublicReportHTTPHandlerName).(http.Handler),
						ctn.Get(constants.AuthHTTPHandlerName).(http.Handler),
						ctn.Get(constants.DashboardHTTPHandlerName).(http.Handler),
						ctn.Get(constants.AccountHTTPHandlerName).(http.Handler),
						ctn.Get(constants.RoleHTTPHandlerName).(http.Handler),
						ctn.Get(constants.PermissionHTTPHandlerName).(http.Handler),
						ctn.Get(constants.GameHTTPHandlerName).(http.Handler),
						ctn.Get(constants.ReportHTTPHandlerName).(http.Handler),
						ctn.Get(constants.OrganizationHTTPHandlerName).(http.Handler),
						ctn.Get(constants.DictionaryHTTPHandlerName).(http.Handler),
						ctn.Get(constants.CurrencyHTTPHandlerName).(http.Handler),
						ctn.Get(constants.WagerSetHTTPHandlerName).(http.Handler),
						ctn.Get(constants.CurrencySetHTTPHandlerName).(http.Handler),
						ctn.Get(constants.FileSetHTTPHandlerName).(http.Handler),
						ctn.Get(constants.LobbyHTTPHandlerName).(http.Handler),
						ctn.Get(constants.ClientInfoHTTPHandlerName).(http.Handler),
					}

					return http.New(ctx, wg, cfg.HTTPConfig, handlers), nil
				},
				Close: func(obj interface{}) error {
					if err := obj.(*http.Server).Shutdown(); err != nil {
						zap.S().Errorf("Error stopping server: %s", err)

						return err
					}

					return nil
				},
			},
			{
				Name: constants.AuthorizerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return jwt.NewProvider(cfg.JWTConfig), nil
				},
			},
			{
				Name: constants.TOTPName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return totp.NewTOTP(cfg.TOTPConfig), nil
				},
			},
			{
				Name: constants.MailgunName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return mailgun.New(cfg.MailgunConfig), nil
				},
			},
			{
				Name: constants.TracerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return tracer.NewTracer(cfg.TracerConfig)
				},
			},

			{
				Name: constants.RPCName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)
					organizationService := ctn.Get(constants.OrganizationServiceName).(*services.OrganizationService)
					gameService := ctn.Get(constants.GameServiceName).(*services.GameService)
					currencyService := ctn.Get(constants.CurrencyServiceName).(*services.CurrencyService)

					return rpc.NewHandler(cfg.RPCConfig, organizationService, gameService, currencyService), nil
				},
			},
			{
				Name: constants.OverlordClientName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return overlord.NewClient(cfg.OverlordConfig)
				},
			},
			{
				Name: constants.ExchangeName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return exchange.NewClient(cfg.ExchangeConfig)
				},
			},
		}

		defs = append(defs, BuildRepositories()...)
		defs = append(defs, BuildServices()...)
		defs = append(defs, BuildHandlers()...)

		if err := builder.Add(defs...); err != nil {
			panic(err)
		}

		container = builder.Build()
	})

	return container
}
