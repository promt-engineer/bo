package config

import (
	"fmt"
	"sync"

	"backoffice/internal/services"
	"backoffice/internal/transport/http"
	"backoffice/internal/transport/queue"
	"backoffice/internal/transport/rpc"
	"backoffice/pkg/auth/jwt"
	"backoffice/pkg/exchange"
	"backoffice/pkg/file"
	"backoffice/pkg/history"
	"backoffice/pkg/mailgun"
	"backoffice/pkg/overlord"
	"backoffice/pkg/pgsql"
	"backoffice/pkg/redis"
	"backoffice/pkg/totp"
	"bitbucket.org/play-workspace/gocommon/tracer"
	"github.com/spf13/viper"
)

var (
	err    error
	config *Config
	once   sync.Once
)

type Config struct {
	Env              string
	LogLevel         string
	FrontURL         string
	SendEmail        string
	ResetPasswordURL string

	PgSQLConfig      *pgsql.Config
	HTTPConfig       *http.Config
	RedisConfig      *redis.Config
	QueueConfig      *queue.Config
	MailgunConfig    *mailgun.Config
	JWTConfig        *jwt.Config
	TOTPConfig       *totp.Config
	TracerConfig     *tracer.Config
	RPCConfig        *rpc.Config
	HistoryConfig    *history.Config
	FileConfig       *file.Config
	LobbyConfig      *services.LobbyConfig
	OverlordConfig   *overlord.Config
	ExchangeConfig   *exchange.Config
	ClientInfoConfig *services.ClientInfoConfig
}

func New() (*Config, error) {
	once.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		httpConfig := viper.Sub("server")
		databaseConfig := viper.Sub("database")
		redisConfig := viper.Sub("redis")
		queueConfig := viper.Sub("queue")
		jwtConfig := viper.Sub("jwt")
		totpConfig := viper.Sub("totp")
		mailgunConfig := viper.Sub("mailgun")
		tracerConfig := viper.Sub("tracer")
		rpcConfig := viper.Sub("rpc")
		historyConfig := viper.Sub("history")
		fileConfig := viper.Sub("file")
		lobbyConfig := viper.Sub("lobby")
		overlordConfig := viper.Sub("overlord")
		exchangeConfig := viper.Sub("exchange")
		clientInfoConfig := viper.Sub("client")

		config.Env = viper.Get("env").(string)
		config.LogLevel = viper.Get("logLevel").(string)
		config.FrontURL = viper.Get("frontURL").(string)
		config.SendEmail = viper.Get("sendEmail").(string)
		config.ResetPasswordURL = viper.Get("resetPasswordURL").(string)

		if err = parseSubConfig(databaseConfig, &config.PgSQLConfig); err != nil {
			return
		}

		if err = parseSubConfig(httpConfig, &config.HTTPConfig); err != nil {
			return
		}

		if err = parseSubConfig(redisConfig, &config.RedisConfig); err != nil {
			return
		}

		if err = parseSubConfig(queueConfig, &config.QueueConfig); err != nil {
			return
		}

		if err = parseSubConfig(mailgunConfig, &config.MailgunConfig); err != nil {
			return
		}

		if err = parseSubConfig(jwtConfig, &config.JWTConfig); err != nil {
			return
		}

		if err = parseSubConfig(totpConfig, &config.TOTPConfig); err != nil {
			return
		}

		if tracerConfig != nil {
			if err := tracerConfig.Unmarshal(&config.TracerConfig); err != nil {
				panic(err)
			}
		} else {
			config.TracerConfig = &tracer.Config{Disabled: true}
		}

		if err = parseSubConfig(rpcConfig, &config.RPCConfig); err != nil {
			return
		}

		if err = parseSubConfig(historyConfig, &config.HistoryConfig); err != nil {
			return
		}

		if err = parseSubConfig(fileConfig, &config.FileConfig); err != nil {
			return
		}

		if err = parseSubConfig(lobbyConfig, &config.LobbyConfig); err != nil {
			return
		}

		if err = parseSubConfig(overlordConfig, &config.OverlordConfig); err != nil {
			return
		}

		if err = parseSubConfig(exchangeConfig, &config.ExchangeConfig); err != nil {
			return
		}

		if err = parseSubConfig(clientInfoConfig, &config.ClientInfoConfig); err != nil {
			return
		}

	})

	return config, err
}

func parseSubConfig[T any](subConfig *viper.Viper, parseTo *T) error {
	if subConfig == nil {
		return fmt.Errorf("can not read %T config: subconfig is nil", parseTo)
	}

	if err := subConfig.Unmarshal(parseTo); err != nil {
		return err
	}

	return nil
}
