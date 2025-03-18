package container

import (
	"backoffice/internal/constants"
	"backoffice/internal/services"
	httpHandlers "backoffice/internal/transport/http/handlers"
	queueHandlers "backoffice/internal/transport/queue/handlers"
	"backoffice/pkg/auth"

	"github.com/sarulabs/di"
)

func BuildHandlers() []di.Def {
	return []di.Def{
		{
			Name: constants.MetaHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				debugService := ctn.Get(constants.DebugServiceName).(*services.DebugService)

				return httpHandlers.NewMetaHandler(debugService), nil
			},
		},
		{
			Name: constants.AuthHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				authz := ctn.Get(constants.AuthorizerName).(auth.Authorizer)
				authService := ctn.Get(constants.AuthenticationServiceName).(*services.AuthenticationService)
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)
				sessionService := ctn.Get(constants.SessionServiceName).(*services.SessionService)

				return httpHandlers.NewAuthHandler(authz, authService, accountService, sessionService), nil
			},
		},
		{
			Name: constants.DashboardHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				return httpHandlers.NewDashboardHandler(), nil
			},
		},
		{
			Name: constants.AccountHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				accountService := ctn.Get(constants.AccountServiceName).(*services.AccountService)
				authorizationService := ctn.Get(constants.AuthorizationServiceName).(*services.AuthorizationService)
				organizationService := ctn.Get(constants.OrganizationServiceName).(*services.OrganizationService)

				return httpHandlers.NewAccountHandler(accountService, authorizationService, organizationService), nil
			},
		},
		{
			Name: constants.RoleHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				authorizationService := ctn.Get(constants.AuthorizationServiceName).(*services.AuthorizationService)
				authenticationService := ctn.Get(constants.AuthenticationServiceName).(*services.AuthenticationService)

				return httpHandlers.NewRoleHandler(authorizationService, authenticationService), nil
			},
		},
		{
			Name: constants.PermissionHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				authorizationService := ctn.Get(constants.AuthorizationServiceName).(*services.AuthorizationService)

				return httpHandlers.NewPermissionHandler(authorizationService), nil
			},
		},
		{
			Name: constants.GameHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				gameService := ctn.Get(constants.GameServiceName).(*services.GameService)
				cfgSender := ctn.Get(constants.ConfigSenderServiceName).(*services.ConfigSenderService)

				return httpHandlers.NewGameHTTPHandler(gameService, cfgSender), nil
			},
		},
		{
			Name: constants.ReportHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				spinService := ctn.Get(constants.SpinServiceName).(*services.SpinService)
				fileService := ctn.Get(constants.FileDownloadingServiceName).(*services.FileDownloadingService)

				return httpHandlers.NewReportHandler(spinService, fileService), nil
			},
		},
		{
			Name: constants.OrganizationHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				organizationService := ctn.Get(constants.OrganizationServiceName).(*services.OrganizationService)
				cfgSender := ctn.Get(constants.ConfigSenderServiceName).(*services.ConfigSenderService)

				return httpHandlers.NewOrganizationHandler(organizationService, cfgSender), nil
			},
		},
		{
			Name: constants.DictionaryHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				gameService := ctn.Get(constants.GameServiceName).(*services.GameService)
				spinService := ctn.Get(constants.SpinServiceName).(*services.SpinService)
				currencyService := ctn.Get(constants.CurrencyServiceName).(*services.CurrencyService)
				organizationService := ctn.Get(constants.OrganizationServiceName).(*services.OrganizationService)

				return httpHandlers.NewDictionaryHandler(gameService, spinService, currencyService, organizationService), nil
			},
		},
		{
			Name: constants.CurrencyHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				currencyService := ctn.Get(constants.CurrencyServiceName).(*services.CurrencyService)
				cfgSender := ctn.Get(constants.ConfigSenderServiceName).(*services.ConfigSenderService)
				fileService := ctn.Get(constants.FileDownloadingServiceName).(*services.FileDownloadingService)

				return httpHandlers.NewCurrencyHandler(currencyService, cfgSender, fileService), nil
			},
		},
		{
			Name: constants.WagerSetHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				wagerSetService := ctn.Get(constants.WagerSetServiceName).(*services.WagerSetService)

				return httpHandlers.NewWagerSetHandler(wagerSetService), nil
			},
		},
		{
			Name: constants.CurrencySetHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				currencySetService := ctn.Get(constants.CurrencySetServiceName).(*services.CurrencySetService)

				return httpHandlers.NewCurrencySetHandler(currencySetService), nil
			},
		},
		{
			Name: constants.CurrencyQueueHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				configSenderService := ctn.Get(constants.ConfigSenderServiceName).(*services.ConfigSenderService)

				return queueHandlers.NewCurrencyHandler(configSenderService), nil
			},
		},
		{
			Name: constants.FileSetHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				fileService := ctn.Get(constants.FileDownloadingServiceName).(*services.FileDownloadingService)

				return httpHandlers.NewFileHandler(fileService), nil
			},
		},
		{
			Name: constants.PublicReportHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				spinService := ctn.Get(constants.SpinServiceName).(*services.SpinService)

				return httpHandlers.NewPublicReportHandler(spinService), nil
			},
		},
		{
			Name: constants.LobbyHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				lobbyService := ctn.Get(constants.LobbyServiceName).(*services.LobbyService)

				return httpHandlers.NewLobbyHTTPHandler(lobbyService), nil
			},
		},
		{
			Name: constants.ClientInfoHTTPHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				clientInfoService := ctn.Get(constants.ClientInfoServiceName).(*services.ClientInfoService)

				return httpHandlers.NewClientInfoHTTPHandler(clientInfoService), nil
			},
		},
	}
}
