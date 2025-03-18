package rpc

import (
	e "backoffice/internal/errors"
	"backoffice/internal/services"
	"backoffice/pkg/backoffice"
	"context"
	"errors"
	"io"

	"go.uber.org/zap"
)

type Handler struct {
	backoffice.BackofficeServer
	cfg                 *Config
	organizationService *services.OrganizationService
	gameService         *services.GameService
	currencyService     *services.CurrencyService
}

func NewHandler(cfg *Config, organizationService *services.OrganizationService, gameService *services.GameService, currencyService *services.CurrencyService) *Handler {
	return &Handler{
		cfg:                 cfg,
		organizationService: organizationService,
		gameService:         gameService,
		currencyService:     currencyService,
	}
}

func (h *Handler) HasAccess(ctx context.Context, in *backoffice.HasAccessIn) (*backoffice.HasAccessOut, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByApiKey(ctx, in.ApiKey)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(e.ErrNotAuthorized)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	if err = h.organizationService.HasAccess(ctx, integrator.ID, in.Game); err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrDoesNotHavePermission) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	return &backoffice.HasAccessOut{HasAccess: true, Integrator: integrator.Name}, nil
}

func (h *Handler) GameList(ctx context.Context, in *backoffice.GameListIn) (*backoffice.GameListOut, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByApiKey(ctx, in.ApiKey)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	games, err := h.gameService.GetIntegratorGameList(ctx, integrator.ID)
	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	return &backoffice.GameListOut{Games: games}, nil
}

func (h *Handler) GameListFull(ctx context.Context, in *backoffice.GameListIn) (*backoffice.GameListOutFull, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByApiKey(ctx, in.ApiKey)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	games, err := h.gameService.GetIntegratorGames(ctx, integrator.ID)
	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	var res = make([]*backoffice.Game, 0, len(games))

	for _, game := range games {

		var aws = make([]*backoffice.WagerSets, 0, len(game.AvailableWagerSets))
		for _, v := range game.AvailableWagerSets {
			newAws := backoffice.WagerSets{Id: v.ID.String(), WagerLevels: v.WagerLevels, DefaultWager: v.DefaultWager}
			aws = append(aws, &newAws)
		}
		var wagerSet backoffice.WagerSets
		if game.WagerSet != nil {
			wagerSet.Id = game.WagerSet.ID.String()
			wagerSet.DefaultWager = game.WagerSet.DefaultWager
			wagerSet.WagerLevels = game.WagerSet.WagerLevels
		}

		res = append(res, &backoffice.Game{
			Id:                  game.ID.String(),
			Name:                game.Name,
			ApiUrl:              game.ApiUrl,
			ClientUrl:           game.ClientUrl,
			IsPublic:            game.IsPublic,
			IsStatisticShown:    game.IsStatisticShown,
			Languages:           game.Languages,
			Currencies:          game.Currencies,
			IsDemo:              game.IsDemo,
			IsFreespins:         game.IsFreespins,
			Rtp:                 game.RTP,
			Volatility:          game.Volatility,
			AvailableRtp:        game.AvailableRTP,
			AvailableVolatility: game.AvailableVolatility,
			UserLocales:         game.UserLocales,
			AvailableWagerSets:  aws,
			GambleDoubleUp:      game.GambleDoubleUp,
			WagerSet:            &wagerSet,
			Provider:            game.OrganizationID.String(),
		})
	}

	return &backoffice.GameListOutFull{Games: res}, nil
}

func (h *Handler) GetProvider(ctx context.Context, in *backoffice.GetProviderIn) (*backoffice.GetProviderOut, error) {
	res, err := h.gameService.GetGameByName(ctx, in.GameName)
	if err != nil {
		return nil, err
	}

	return &backoffice.GetProviderOut{Provider: res.Organization.Name}, err
}

func (h *Handler) HealthCheck(stream backoffice.Backoffice_HealthCheckServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if err = stream.Send(msg); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}
}

func (h *Handler) GetIntegratorGameSettings(ctx context.Context, in *backoffice.IntegratorGameSettingsIn) (*backoffice.IntegratorGameSettingsOut, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByApiKey(ctx, in.ApiKey)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(e.ErrNotAuthorized)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	integratorGameSettings, err := h.organizationService.GetIntegratorGameSettings(ctx, integrator.ID, in.Game, in.Currency)

	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	var WagerLevels []int64
	if integratorGameSettings.WagerSet == nil {
		WagerLevels = nil
	} else {
		WagerLevels = integratorGameSettings.WagerSet.WagerLevels
	}

	return &backoffice.IntegratorGameSettingsOut{Wagers: WagerLevels, Rtp: integratorGameSettings.RTP, Volatility: integratorGameSettings.Volatility, ShortLink: integratorGameSettings.ShortLink}, nil
}

func (h *Handler) GetGameData(ctx context.Context, in *backoffice.GameDataByIntegratorName) (*backoffice.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByName(ctx, in.Name)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	game, err := h.gameService.GetIntegratorGame(ctx, integrator.ID, in.Game)
	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	var aws = make([]*backoffice.WagerSets, 0, len(game.AvailableWagerSets))
	for _, v := range game.AvailableWagerSets {
		newAws := backoffice.WagerSets{Id: v.ID.String(), WagerLevels: v.WagerLevels, DefaultWager: v.DefaultWager}
		aws = append(aws, &newAws)
	}

	return &backoffice.Game{
		Id:                  game.ID.String(),
		Name:                game.Name,
		ApiUrl:              game.ApiUrl,
		ClientUrl:           game.ClientUrl,
		IsPublic:            game.IsPublic,
		IsStatisticShown:    game.IsStatisticShown,
		Languages:           game.Languages,
		Currencies:          game.Currencies,
		IsDemo:              game.IsDemo,
		IsFreespins:         game.IsFreespins,
		Rtp:                 game.RTP,
		Volatility:          game.Volatility,
		AvailableRtp:        game.AvailableRTP,
		AvailableVolatility: game.AvailableVolatility,
		UserLocales:         game.UserLocales,
		AvailableWagerSets:  aws,
	}, nil
}

func (h *Handler) GetGameDataByApi(ctx context.Context, in *backoffice.HasAccessIn) (*backoffice.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByApiKey(ctx, in.ApiKey)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	game, err := h.gameService.GetIntegratorGame(ctx, integrator.ID, in.Game)
	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	var aws = make([]*backoffice.WagerSets, 0, len(game.AvailableWagerSets))
	for _, v := range game.AvailableWagerSets {
		newAws := backoffice.WagerSets{Id: v.ID.String(), WagerLevels: v.WagerLevels, DefaultWager: v.DefaultWager}
		aws = append(aws, &newAws)
	}

	return &backoffice.Game{
		Id:                  game.ID.String(),
		Name:                game.Name,
		ApiUrl:              game.ApiUrl,
		ClientUrl:           game.ClientUrl,
		IsPublic:            game.IsPublic,
		IsStatisticShown:    game.IsStatisticShown,
		Languages:           game.Languages,
		Currencies:          game.Currencies,
		IsDemo:              game.IsDemo,
		IsFreespins:         game.IsFreespins,
		Rtp:                 game.RTP,
		Volatility:          game.Volatility,
		AvailableRtp:        game.AvailableRTP,
		AvailableVolatility: game.AvailableVolatility,
		UserLocales:         game.UserLocales,
		AvailableWagerSets:  aws,
	}, nil
}

func (h *Handler) GetCurrencies(ctx context.Context, in *backoffice.CurrenciesIn) (*backoffice.CurrenciesOut, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	filters := make(map[string]interface{})
	for _, filter := range in.CurrenciesFilter {
		filters[filter.Key] = filter.Value
	}

	currencies, err := h.currencyService.CurrencyGetAll(ctx, filters)

	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	var result backoffice.CurrenciesOut
	for _, currency := range currencies {
		result.Currencies = append(result.Currencies, &backoffice.Currency{
			Currency:     currency.Alias,
			BaseCurrency: currency.BaseCurrency,
		})
	}

	return &result, nil
}

func (h *Handler) GetMultiplierByCurrency(ctx context.Context, in *backoffice.GetMultiplierIn) (*backoffice.GetMultiplierOut, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByApiKey(ctx, in.ApiKey)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	organizationsPairIDs, err := h.organizationService.GetOrganizationPairsByIntegrator(ctx, integrator.ID)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(err)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	currencyMultipliers, err := h.currencyService.GetCurrencyMultipliersByOrganizationPairs(ctx, organizationsPairIDs, in.Currency)

	if err != nil {
		zap.S().Error(err)

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	var providerMultipliers []*backoffice.ProviderMultiplierEntry
	for _, multiplier := range currencyMultipliers {
		provider := multiplier.ProviderIntegratorPair.ProviderID.String()
		entry := &backoffice.ProviderMultiplierEntry{
			Provider:   provider,
			Multiplier: multiplier.Multiplier,
		}
		providerMultipliers = append(providerMultipliers, entry)
	}

	return &backoffice.GetMultiplierOut{
		ProviderMultipliers: providerMultipliers,
	}, nil
}

func (h *Handler) GetIntegratorApiKey(ctx context.Context, in *backoffice.IntegratorApiKeyIn) (*backoffice.IntegratorApiKeyOut, error) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.MaxProcessingTime)
	defer cancel()

	integrator, err := h.organizationService.GetByName(ctx, in.Integrator)
	if err != nil {
		zap.S().Error(err)

		if errors.Is(err, e.ErrEntityNotFound) {
			return nil, WrapInGRPCError(e.ErrNotAuthorized)
		}

		return nil, WrapInGRPCError(e.ErrInternal)
	}

	return &backoffice.IntegratorApiKeyOut{ApiKey: integrator.ApiKey}, nil
}
