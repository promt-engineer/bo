package services

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"backoffice/internal/transport/http/requests"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type GameService struct {
	gameRepo         repositories.GameRepository
	organizationRepo repositories.OrganizationRepository
}

func NewGameService(gameRepo repositories.GameRepository, organizationRepo repositories.OrganizationRepository) *GameService {
	return &GameService{gameRepo: gameRepo, organizationRepo: organizationRepo}
}

/*
gameDivisors is a map that defines the divisor values for specific games.

	These values are used to validate the wager levels for each game.
*/
var gameDivisors = map[string]int{
	"admiral-wilds":            90,
	"frozen-fruits-flexiways":  30,
	"lucky-skulls-bonanza":     200,
	"irish-riches-bonanza":     200,
	"wild-dragon-respin":       100,
	"coral-reef-flexiways":     200,
	"lucky-santa-bonanza":      200,
	"fortune-777-respin":       100,
	"cleos-riches-flexiways":   200,
	"sweet-mystery-flexiways":  100,
	"quest-of-ra":              200,
	"vampire-vault-hold-n-win": 200,
	"yakuza-clash-hold-n-win":  200,
	"double-fortune-panda":     200,
	"squid-gold-x2":            200,
	"brazilian-mask-fire":      200,
}

func (s *GameService) Create(ctx context.Context, req *requests.GameRequest) (*entities.Game, error) {
	_, err := s.gameRepo.GetBy(ctx, map[string]interface{}{"name": req.Name})
	if err == nil {
		return nil, e.ErrEntityAlreadyExist
	}

	if err != nil && !errors.Is(err, e.ErrEntityNotFound) {
		return nil, err
	}

	organization, err := s.organizationRepo.Get(ctx, map[string]interface{}{"id": req.OrganizationID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotProvider
	}

	divisor, exists := gameDivisors[req.Name]
	if exists {
		if err := s.validateWagerSet(ctx, req.Name, req.WagerSetID, divisor); err != nil {
			return nil, fmt.Errorf("invalid wager set values for game %s", req.Name)
		}
	}
	game := &entities.Game{
		ID:                   uuid.New(),
		OrganizationID:       req.OrganizationID,
		Name:                 req.Name,
		Jurisdictions:        req.Jurisdictions,
		Currencies:           req.Currencies,
		Languages:            req.Languages,
		UserLocales:          req.UserLocales,
		ApiUrl:               req.ApiURL,
		ClientUrl:            req.ClientURL,
		WagerSetID:           req.WagerSetID,
		IsPublic:             *req.IsPublic,
		IsStatisticShown:     *req.IsStatisticShown,
		IsDemo:               *req.IsDemo,
		IsFreespins:          *req.IsFreeSpins,
		RTP:                  req.RTP,
		Volatility:           req.Volatility,
		AvailableRTP:         req.AvailableRTP,
		AvailableVolatility:  req.AvailableVolatility,
		OnlineVolatility:     *req.OnlineVolatility,
		AvailableWagerSetsID: req.AvailableWagerSetsID,
		GambleDoubleUp:       req.GambleDoubleUp,
	}

	g, err := s.gameRepo.Create(ctx, game)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *GameService) GetGame(ctx context.Context, gameID uuid.UUID) (*entities.Game, error) {
	return s.gameRepo.GetBy(ctx, map[string]interface{}{"id": gameID})
}

func (s *GameService) GetGameByName(ctx context.Context, name string) (*entities.Game, error) {
	return s.gameRepo.GetBy(ctx, map[string]interface{}{"name": name})
}

func (s *GameService) AllForStat(ctx context.Context, organizationID *uuid.UUID) ([]*entities.Game, error) {
	return s.gameRepo.All(ctx, organizationID, map[string]interface{}{"is_statistic_shown": true})
}

func (s *GameService) AllPublic(ctx context.Context, organizationID *uuid.UUID) ([]*entities.Game, error) {
	return s.gameRepo.All(ctx, organizationID, map[string]interface{}{"is_public": true})
}

func (s *GameService) GetDictionaries(ctx context.Context, organizationID *uuid.UUID, dictType string) ([]string, error) {
	return s.gameRepo.GetDictionaries(ctx, organizationID, dictType)
}

func (s *GameService) IDs(ctx context.Context, organizationID *uuid.UUID) ([]uuid.UUID, error) {
	games, err := s.AllPublic(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	var list []uuid.UUID

	for _, game := range games {
		list = append(list, game.ID)
	}

	return list, nil
}

func (s *GameService) IDsString(ctx context.Context, organizationID *uuid.UUID) ([]string, error) {
	list, err := s.IDs(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return lo.Map(list, func(item uuid.UUID, index int) string {
		return item.String()
	}), nil
}

func (s *GameService) GetIntegratorGameNames(ctx context.Context, organizationID uuid.UUID) ([]string, error) {
	//games, err := s.gameRepo.GetOrganizationGameList(ctx, organizationID)
	games, err := s.gameRepo.GetIntegratorGameList(ctx, organizationID)
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	list := make([]string, 0, len(games))

	for _, game := range games {
		list = append(list, game.Name)
	}

	return list, nil
}

func (s *GameService) GetProviderGameIDs(ctx context.Context, organizationID uuid.UUID) ([]uuid.UUID, error) {
	games, err := s.gameRepo.GetOrganizationGameList(ctx, organizationID)
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	list := make([]uuid.UUID, 0, len(games))

	for _, game := range games {
		list = append(list, game.ID)
	}

	return list, nil
}

func (s *GameService) GetIntegratorGameList(ctx context.Context, organizationID uuid.UUID) ([]string, error) {
	games, err := s.gameRepo.GetIntegratorGameList(ctx, organizationID)
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	list := make([]string, 0, len(games))

	for _, game := range games {
		list = append(list, game.Name)
	}

	return list, nil
}

func (s *GameService) GetIntegratorGames(ctx context.Context, organizationID uuid.UUID) ([]*entities.Game, error) {
	games, err := s.gameRepo.GetIntegratorGameList(ctx, organizationID)
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	return games, nil
}

func (s *GameService) GetGameList(ctx context.Context, filter map[string]interface{}) ([]*entities.Game, error) {
	return s.gameRepo.GetAllByFilter(ctx, filter)
}

func (s *GameService) Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) ([]*entities.Game, int64, error) {
	return s.gameRepo.Paginate(ctx, filters, order, limit, offset)
}

func (s *GameService) Update(ctx context.Context, gameID uuid.UUID, req *requests.GameRequest) (*entities.Game, error) {
	organization, err := s.organizationRepo.Get(ctx, map[string]interface{}{"id": req.OrganizationID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotProvider
	}

	divisor, exists := gameDivisors[req.Name]
	if exists {
		if err := s.validateWagerSet(ctx, req.Name, req.WagerSetID, divisor); err != nil {
			return nil, fmt.Errorf("invalid wager set values for game %s", req.Name)
		}
	}

	gameParam := map[string]interface{}{
		"organization_id":         req.OrganizationID,
		"name":                    req.Name,
		"jurisdictions":           req.Jurisdictions,
		"currencies":              req.Currencies,
		"languages":               req.Languages,
		"user_locales":            req.UserLocales,
		"api_url":                 req.ApiURL,
		"client_url":              req.ClientURL,
		"wager_set_id":            req.WagerSetID,
		"is_public":               *req.IsPublic,
		"is_statistic_shown":      *req.IsStatisticShown,
		"is_demo":                 *req.IsDemo,
		"is_freespins":            *req.IsFreeSpins,
		"rtp":                     req.RTP,
		"volatility":              req.Volatility,
		"available_rtp":           req.AvailableRTP,
		"available_volatility":    req.AvailableVolatility,
		"online_volatility":       *req.OnlineVolatility,
		"available_wager_sets_id": req.AvailableWagerSetsID,
		"gamble_double_up":        req.GambleDoubleUp,
	}

	g, err := s.gameRepo.Update(ctx, gameID, gameParam)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *GameService) Delete(ctx context.Context, gameID uuid.UUID) error {
	game, err := s.gameRepo.GetBy(ctx, map[string]interface{}{"id": gameID})
	if err != nil {
		return err
	}

	err = s.gameRepo.Delete(ctx, game)
	if err != nil {
		return err
	}

	return nil
}

func (s *GameService) AddValueToTheDictionary(ctx context.Context, organizationID *uuid.UUID, dictType, value string) (string, error) {
	availableValues, err := s.gameRepo.GetDictionaries(ctx, organizationID, dictType)
	if err != nil {
		return "", err
	}

	for _, v := range availableValues {
		if v == value {
			return "", fmt.Errorf("such record already exists")
		}
	}

	return s.gameRepo.AddValueToTheDictionary(ctx, organizationID, dictType, value)
}

func (s *GameService) DeleteValueFromTheDictionary(ctx context.Context, organizationID *uuid.UUID, dictType, value string) error {
	availableValues, err := s.gameRepo.GetDictionaries(ctx, organizationID, dictType)
	if err != nil {
		return err
	}

	valueIsAbsent := true
	for _, v := range availableValues {
		if v == value {
			valueIsAbsent = false
			break
		}
	}

	if valueIsAbsent {
		return fmt.Errorf("no such record exists")
	}

	return s.gameRepo.RemoveValueFromDictionary(ctx, organizationID, dictType, value)
}

func (s *GameService) GetIntegratorGame(ctx context.Context, organizationID uuid.UUID, gameName string) (*entities.Game, error) {
	game, err := s.gameRepo.GetIntegratorGame(ctx, organizationID, gameName)
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	return game, nil
}

func (s *GameService) GetAvailableWagerSetsByIDs(ctx context.Context, game *entities.Game) (wagerSets []entities.WagerSet, err error) {
	return s.gameRepo.GetAvailableWagerSetsByIDs(ctx, game)
}

func (s *GameService) validateWagerSet(ctx context.Context, gameName string, wagerSetID uuid.UUID, divisor int) error {
	wagerSet, err := s.gameRepo.GetWagerSetByID(ctx, wagerSetID)
	if err != nil {
		return err
	}

	for _, value := range wagerSet.WagerLevels {
		if value%int64(divisor) != 0 {
			return fmt.Errorf("invalid wager set values for game %s", gameName)
		}
	}

	return nil
}
