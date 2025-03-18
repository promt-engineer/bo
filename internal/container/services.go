package container

import (
	"backoffice/internal/config"
	"backoffice/internal/constants"
	"backoffice/internal/entities"
	"backoffice/internal/repositories"
	"backoffice/internal/services"
	"backoffice/internal/transport/queue"
	"backoffice/pkg/auth"
	"backoffice/pkg/exchange"
	"backoffice/pkg/history"
	"backoffice/pkg/mailgun"
	"backoffice/pkg/overlord"

	"github.com/sarulabs/di"
)

func BuildServices() []di.Def {
	return []di.Def{
		{
			Name: constants.AccountServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				repo := ctn.Get(constants.AccountRepositoryName).(repositories.AccountRepository)
				sessionService := ctn.Get(constants.SessionServiceName).(*services.SessionService)
				mailingService := ctn.Get(constants.MailingServiceName).(*services.MailingService)

				return services.NewAccountService(repo, sessionService, mailingService), nil
			},
		},
		{
			Name: constants.AuthenticationServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				authz := ctn.Get(constants.AuthorizerName).(auth.Authorizer)
				authorizationService := ctn.Get(constants.AuthorizationServiceName).(*services.AuthorizationService)
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)
				sessionService := ctn.Get(constants.SessionServiceName).(*services.SessionService)
				refreshTokenRepo := ctn.Get(constants.RefreshTokenRepositoryName).(repositories.TokenRepository)
				organizationService := ctn.Get(constants.OrganizationServiceName).(*services.OrganizationService)

				return services.NewAuthenticationService(authz, authorizationService, accountService, sessionService, organizationService, refreshTokenRepo), nil
			},
		},
		{
			Name: constants.SessionServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				repo := ctn.Get(constants.SessionRepositoryName).(repositories.SessionRepository)

				return services.NewSessionService(repo), nil
			},
		},
		{
			Name: constants.AuthorizationServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)
				roleRepo := ctn.Get(constants.RoleRepositoryName).(repositories.RoleRepository)
				permRepo := ctn.Get(constants.PermissionRepositoryName).(repositories.PermissionRepository)

				return services.NewAuthorizationService(accountService, roleRepo, permRepo), nil
			},
		},
		{
			Name: constants.GameServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				gameRepo := ctn.Get(constants.GameRepositoryName).(repositories.GameRepository)
				organizationRepo := ctn.Get(constants.OrganizationRepositoryName).(repositories.OrganizationRepository)

				return services.NewGameService(gameRepo, organizationRepo), nil
			},
		},
		{
			Name: constants.SpinServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				gameService := ctn.Get(constants.GameServiceName).(*services.GameService)
				historyClient := ctn.Get(constants.HistoryName).(history.Client)

				return services.NewSpinService(gameService, historyClient), nil
			},
		},
		{
			Name: constants.QueueName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return queue.NewQueue(cfg.QueueConfig), nil
			},
		},
		{
			Name: constants.OrganizationServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				repo := ctn.Get(constants.OrganizationRepositoryName).(repositories.OrganizationRepository)
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)
				gameService := ctn.Get(constants.GameServiceName).(*services.GameService)

				return services.NewOrganizationService(repo, accountService, gameService), nil
			},
		},
		{
			Name: constants.CurrencyServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				repo := ctn.Get(constants.CurrencyRepositoryName).(repositories.CurrencyRepository)
				organizationService := ctn.Get(constants.OrganizationServiceName).(*services.OrganizationService)
				exchangeClient := ctn.Get(constants.ExchangeName).(exchange.Client)

				return services.NewCurrencyService(repo, organizationService, exchangeClient), nil
			},
		},
		{
			Name: constants.LobbyServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				gameService := ctn.Get(constants.GameServiceName).(*services.GameService)
				overlordClient := ctn.Get(constants.OverlordClientName).(overlord.Client)

				return services.NewLobbyService(cfg.LobbyConfig, overlordClient, gameService), nil
			}},
		{
			Name: constants.MailingServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				mailgunClient := ctn.Get(constants.MailgunName).(*mailgun.Client)
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return services.NewMailingService(mailgunClient, cfg.FrontURL, cfg.SendEmail, cfg.ResetPasswordURL), nil
			},
		},
		{
			Name: constants.WagerSetServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				wagerSetRepo := ctn.Get(constants.WagerSetRepositoryName).(repositories.BaseRepository[entities.WagerSet])

				return services.NewWagerSetService(wagerSetRepo), nil
			},
		},
		{
			Name: constants.CurrencySetServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				currencyService := ctn.Get(constants.CurrencyServiceName).(*services.CurrencyService)
				currencySetRepo := ctn.Get(constants.CurrencySetRepositoryName).(repositories.BaseRepository[entities.CurrencySet])

				return services.NewCurrencySetService(currencyService, currencySetRepo), nil
			},
		},
		{
			Name: constants.ConfigSenderServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				publisher := ctn.Get(constants.QueueName).(*queue.Queue)
				currencyService := ctn.Get(constants.CurrencyServiceName).(*services.CurrencyService)
				gameService := ctn.Get(constants.GameServiceName).(*services.GameService)

				return services.NewConfigSenderService(publisher, currencyService, gameService), nil
			},
		},
		{
			Name: constants.DebugServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				debugRepo := ctn.Get(constants.DebugRepositoryName).(repositories.DebugRepository)

				return services.NewDebugService(debugRepo), nil
			},
		},
		{
			Name: constants.FileDownloadingServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)
				spinService := ctn.Get(constants.SpinServiceName).(*services.SpinService)
				fileRepo := ctn.Get(constants.FileRepositoryName).(repositories.FileRepository)

				return services.NewFileDownloadingService(cfg.FileConfig, cfg.ClientInfoConfig, spinService, fileRepo), nil
			},
		},
		{
			Name: constants.ClientInfoServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return services.NewClientInfoService(cfg.ClientInfoConfig), nil
			}},
	}
}
