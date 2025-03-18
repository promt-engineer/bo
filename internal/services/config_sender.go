package services

import (
	"backoffice/internal/constants"
	"backoffice/internal/entities"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Sender interface {
	Send(ctx context.Context, queueName string, msgType string, payload interface{}) error
}

type Notifier interface {
	Notify(ev any)
}

type ConfigSenderService struct {
	sender Sender

	gameService     *GameService
	currencyService *CurrencyService
}

func NewConfigSenderService(sender Sender, currencyService *CurrencyService, gameService *GameService) *ConfigSenderService {
	return &ConfigSenderService{sender: sender, currencyService: currencyService, gameService: gameService}
}

func (s *ConfigSenderService) Notify(ev any) {
	if err := s.SendCurrencyToLord(context.Background()); err != nil {
		zap.S().Error(err)
	}
}

func (s *ConfigSenderService) SendCurrencyToLord(ctx context.Context) error {
	zap.S().Info("get notify signal for prepare data to lord")
	muls, err := s.currencyService.All(ctx)
	if err != nil {
		return err
	}

	toSend := make([]*entities.CurrencyGameConfig, 0)
	games, err := s.gameService.AllPublic(ctx, nil)
	if err != nil {
		return err
	}

	gamesMap := map[uuid.UUID][]*entities.Game{}

	lo.ForEach(games, func(item *entities.Game, index int) {
		gamesMap[item.OrganizationID] = append(gamesMap[item.OrganizationID], item)
	})

	mulsReady := entities.GroupCurrencyMultiplier(muls)

	lo.ForEach(mulsReady, func(item *entities.GroupedCurrencyMultiplier, index int) {
		if item.ProviderIntegratorPair != nil {
			providerGames := gamesMap[item.ProviderIntegratorPair.ProviderID]
			lo.ForEach(providerGames, func(game *entities.Game, index int) {
				toSend = append(toSend,
					&entities.CurrencyGameConfig{
						IntegratorName: item.ProviderIntegratorPair.Integrator.Name,
						ProviderName:   item.ProviderIntegratorPair.Provider.Name,
						GameName:       game.Name,
						GameID:         game.ID,

						DefaultWager:        game.WagerSet.DefaultWager,
						WagerLevels:         game.WagerSet.WagerLevels,
						Multipliers:         item.Multipliers,
						AvailableRTP:        game.AvailableRTP,
						AvailableVolatility: game.AvailableVolatility,
						OnlineVolatility:    game.OnlineVolatility,
						GambleDoubleUp:      game.GambleDoubleUp,
						Synonyms:            item.Synonyms,
					})
			})
		}
	})

	if len(toSend) == 0 {
		return fmt.Errorf("currency config is empty, nothing to send to lord")
	}

	zap.S().Info("prepare data for lord - DONE")

	return s.sender.Send(ctx, constants.QueueOverlordName, constants.MsgCurrencyConfigType, toSend)
}
