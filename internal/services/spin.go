package services

import (
	"backoffice/internal/entities"
	"backoffice/pkg/history"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type SpinService struct {
	gameService   *GameService
	historyClient history.Client
}

func NewSpinService(gameService *GameService, historyClient history.Client) *SpinService {
	return &SpinService{gameService: gameService, historyClient: historyClient}
}

func (s *SpinService) FinancialReport(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase) (*entities.FinancialReport, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, err
	}

	rep, err := s.historyClient.GetFinancialReport(ctx, fb)
	if err != nil {
		return nil, err
	}

	return entities.FinancialReportFromHistory(rep), nil
}

func (s *SpinService) Paginate(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase, order string, limit int, page int) (
	pagination entities.Pagination[entities.Spin], err error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return pagination, err
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return pagination, err
	}

	out, err := s.historyClient.GetSpins(ctx, &history.GetFinancialIn{
		Page:  uint64(page),
		Limit: uint64(limit),
		Order: order,
		Base:  fb,
	})
	if err != nil {
		return
	}

	lo.ForEach(out.Items, func(item *history.SpinOut, index int) {
		pagination.Items = append(pagination.Items, entities.SpinFromHistory(item))
	})

	pagination.Total = int(out.Total)
	pagination.Limit = int(out.Limit)
	pagination.CurrentPage = int(out.CurrentPage)

	return pagination, nil
}

func (s *SpinService) PaginateGrouped(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase, order string, limit, page int, groupBy []string) (
	pagination entities.Pagination[entities.GroupedSpin], err error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return pagination, err
	}

	out, err := s.historyClient.GetSpins(ctx, &history.GetFinancialIn{
		Page:    uint64(page),
		Limit:   uint64(limit),
		Order:   order,
		GroupBy: groupBy,
		Base:    fb,
	})
	if err != nil {
		return
	}

	lo.ForEach(out.Items, func(item *history.SpinOut, index int) {
		gs := &entities.GroupedSpin{
			CreatedAt:      item.CreatedAt.AsTime(),
			UpdatedAt:      item.UpdatedAt.AsTime(),
			Game:           item.Game,
			Integrator:     item.Integrator,
			Operator:       item.Operator,
			Host:           &item.Host,
			UserID:         &item.InternalUserId,
			ExternalUserID: &item.ExternalUserId,
			Currency:       item.Currency,
			Wager:          float64(item.Wager),
			BaseAward:      float64(item.BaseAward),
			BonusAward:     float64(item.BonusAward),
			FinalAward:     float64(item.FinalAward),
			IsShown:        item.IsShown,
			IsPFR:          item.IsPfr,
		}

		if item.SessionToken != "" {
			gs.SessionToken = uuid.MustParse(item.SessionToken)
		}

		if item.GameId != "" {
			gs.GameID = uuid.MustParse(item.GameId)
		}

		pagination.Items = append(pagination.Items, gs)
	})

	pagination.Total = int(out.Total)
	pagination.Limit = int(out.Limit)
	pagination.CurrentPage = int(out.CurrentPage)

	return pagination, nil
}

func (s *SpinService) PaginateGamingSession(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase, order string, limit, page int) (
	pagination entities.Pagination[entities.GamingSession], err error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return pagination, err
	}

	out, err := s.historyClient.GetSessions(ctx, &history.GetFinancialIn{
		Page:  uint64(page),
		Limit: uint64(limit),
		Order: order,
		Base:  fb,
	})
	if err != nil {
		return
	}

	lo.ForEach(out.Items, func(item *history.GameSessionOut, index int) {
		pagination.Items = append(pagination.Items, entities.GamingSessionFromHistory(item))
	})

	pagination.Total = int(out.Total)
	pagination.Limit = int(out.Limit)
	pagination.CurrentPage = int(out.CurrentPage)

	return pagination, nil
}

func (s *SpinService) Session(ctx context.Context, organizationID *uuid.UUID, id uuid.UUID, currency string) (*entities.GamingSession, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	out, err := s.historyClient.GetSession(ctx, gameIDs, id.String(), currency)
	if err != nil {
		return nil, err
	}

	return entities.GamingSessionFromHistory(out), nil
}

func (s *SpinService) GameSession(ctx context.Context, gameName string, id uuid.UUID, currency string) (*entities.GamingSession, error) {
	game, err := s.gameService.GetGameByName(ctx, gameName)
	if err != nil {
		return nil, err
	}

	var gameIDs []string
	gameIDs = append(gameIDs, game.ID.String())

	out, err := s.historyClient.GetSession(ctx, gameIDs, id.String(), currency)
	if err != nil {
		return nil, err
	}

	return entities.GamingSessionFromHistory(out), nil
}

func (s *SpinService) AllSpins(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase) ([]*entities.Spin, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, err
	}

	out, err := s.historyClient.GetAllSpins(ctx, fb)
	if err != nil {
		return nil, err
	}

	return lo.Map(out.Spins, func(item *history.SpinOut, index int) *entities.Spin {
		return entities.SpinFromHistory(item)
	}), nil
}

func (s *SpinService) AllGamingSessions(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase) ([]*entities.GamingSession, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, err
	}

	out, err := s.historyClient.AllGamingSession(ctx, fb)
	if err != nil {
		return nil, err
	}

	return lo.Map(out.Sessions, func(item *history.GameSessionOut, index int) *entities.GamingSession {
		return entities.GamingSessionFromHistory(item)
	}), nil
}

func (s *SpinService) Currencies(ctx context.Context, organizationID *uuid.UUID, filters *entities.FinancialBase) ([]string, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, err
	}

	return s.historyClient.GetCurrencies(ctx, fb)
}

func (s *SpinService) Hosts(ctx context.Context, organizationID *uuid.UUID) ([]string, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	filters := &entities.FinancialBase{}

	fb, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, err
	}

	return s.historyClient.GetHosts(ctx, fb)
}

func (s *SpinService) IntegratorOperatorsMap(ctx context.Context, organizationID *uuid.UUID) (map[string][]string, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return s.historyClient.IntegratorOperatorsMap(ctx, gameIDs)
}

func (s *SpinService) GetSpin(ctx context.Context, roundID, currency string) (*entities.Spin, error) {
	out, err := s.historyClient.GetSpin(ctx, roundID, currency)
	if err != nil {
		return nil, err
	}

	return entities.SpinFromHistory(out.Item), nil
}

// TODO: make generic
func (s *SpinService) AggregatedReportByGame(ctx context.Context, organizationID *uuid.UUID, currency string, country *string, filters *entities.AggregateFilters) ([]*entities.AggregatedReportByGame, error) {
	allFilter, pfrFilter, err := s.aggregatedReportFilters(ctx, organizationID, filters)
	if err != nil {
		return nil, err
	}

	if country != nil {
		allFilter.Country = country
		pfrFilter.Country = country
	}

	allReps, err := s.historyClient.GetAggregatedReportByGame(ctx, allFilter)
	if err != nil {
		return nil, err
	}

	pfrReps, err := s.historyClient.GetAggregatedReportByGame(ctx, pfrFilter)
	if err != nil {
		return nil, err
	}

	finalReps := map[string]*entities.AggregatedReportByGame{}

	for _, rep := range allReps {
		finalReps[aggregatedReportByGameKey(rep)] = entities.AggregatedReportByGameFromHistory(rep)
	}

	for _, rep := range pfrReps {
		exRep, ok := finalReps[aggregatedReportByGameKey(rep)]
		if !ok {
			exRep = entities.AggregatedReportByGameFromHistory(rep)
		}

		exRep.PFRWager = rep.Wager
		exRep.PFRAward = rep.Award
		exRep.PFRRoundCount = int(rep.RoundCount)

		finalReps[aggregatedReportByGameKey(rep)] = exRep
	}

	final := lo.Values(finalReps)

	lo.ForEach(final, func(item *entities.AggregatedReportByGame, index int) {
		item.Compute()
	})

	return final, nil
}

func (s *SpinService) AggregatedReportByCountry(ctx context.Context, organizationID *uuid.UUID, currency string, game *string, filters *entities.AggregateFilters) ([]*entities.AggregatedReportByCountry, error) {
	allFilter, pfrFilter, err := s.aggregatedReportFilters(ctx, organizationID, filters)
	if err != nil {
		return nil, err
	}

	if game != nil {
		allFilter.Game = game
		pfrFilter.Game = game
	}

	allReps, err := s.historyClient.GetAggregatedReportByCountry(ctx, allFilter)
	if err != nil {
		return nil, err
	}

	pfrReps, err := s.historyClient.GetAggregatedReportByCountry(ctx, pfrFilter)
	if err != nil {
		return nil, err
	}

	finalReps := map[string]*entities.AggregatedReportByCountry{}

	for _, rep := range allReps {
		finalReps[aggregatedReportByCountryKey(rep)] = entities.AggregatedReportByCountryFromHistory(rep)
	}

	for _, rep := range pfrReps {
		exRep, ok := finalReps[aggregatedReportByCountryKey(rep)]
		if !ok {
			exRep = entities.AggregatedReportByCountryFromHistory(rep)
		}

		exRep.PFRWager = rep.GetWager()
		exRep.PFRAward = rep.GetAward()

		finalReps[aggregatedReportByCountryKey(rep)] = exRep
	}

	final := lo.Values(finalReps)

	lo.ForEach(final, func(item *entities.AggregatedReportByCountry, index int) {
		item.Compute()
	})

	return final, nil
}

type AggregatedReportByField interface {
	entities.AggregatedReportByCountry
	entities.AggregatedReportByGame
}

func aggregatedReportByCountryKey(ar *history.GetAggregatedReportByCountryItem) string {
	return fmt.Sprintf("%v/%v", ar.Country, ar.Currency)
}

func aggregatedReportByGameKey(ar *history.GetAggregatedReportByGameItem) string {
	return fmt.Sprintf("%v/%v", ar.GameId, ar.Currency)
}

func (s *SpinService) aggregatedReportFilters(ctx context.Context, organizationID *uuid.UUID, filters *entities.AggregateFilters) (
	*history.GetAggregatedReportFilters, *history.GetAggregatedReportFilters, error) {
	gameIDs, err := s.gameService.IDsString(ctx, organizationID)
	if err != nil {
		return nil, nil, err
	}

	allFilter, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, nil, err
	}

	tr := true
	filters.IsPFR = &tr

	pfrFilter, err := filters.ToHistoryFilters(gameIDs)
	if err != nil {
		return nil, nil, err
	}

	return allFilter, pfrFilter, nil
}

func (s *SpinService) UserReport(ctx context.Context, organizationID *uuid.UUID, id, currency string) (
	report *entities.UserReport, err error) {
	filters := &entities.FinancialBase{}
	filters.ExternalUserID = id
	filters.Currency = &currency

	spins, err := s.AllSpins(ctx, organizationID, filters)
	if err != nil {
		return
	}

	return entities.UserReportFromSpins(spins).Compute(), nil
}
