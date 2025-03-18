package services

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type OrganizationService struct {
	repo           repositories.OrganizationRepository
	accountService *AccountService
	gameService    *GameService
}

func NewOrganizationService(repo repositories.OrganizationRepository, accountService *AccountService, gameService *GameService) *OrganizationService {
	return &OrganizationService{
		repo:           repo,
		accountService: accountService,
		gameService:    gameService,
	}
}

func (s *OrganizationService) Create(ctx context.Context, status uint8, name, t string) (*entities.Organization, error) {
	return s.repo.Create(ctx, &entities.Organization{
		ID:     uuid.New(),
		Name:   name,
		Type:   t,
		ApiKey: uuid.New().String(),
		Status: &status,
	})
}

func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, status uint8, name, t string) (*entities.Organization, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, &entities.Organization{
		ID:     organization.ID,
		Name:   name,
		Type:   t,
		ApiKey: organization.ApiKey,
		Status: &status,
	})
}

func (s *OrganizationService) Delete(ctx context.Context, accountID, organizationID uuid.UUID) error {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": organizationID})
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, accountID, organization)
}

func (s *OrganizationService) Get(ctx context.Context, id uuid.UUID) (*entities.Organization, error) {
	return s.repo.Get(ctx, map[string]interface{}{"id": id})
}

func (s *OrganizationService) Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) ([]*entities.Organization, int64, error) {
	return s.repo.Paginate(ctx, filters, order, limit, offset)
}

func (s *OrganizationService) GetIntegratorsByProvider(ctx context.Context, providerID uuid.UUID) (organizations []*entities.Organization, err error) {
	return s.repo.GetIntegratorsByProvider(ctx, providerID)
}

func (s *OrganizationService) Assign(ctx context.Context, accountID, organizationID uuid.UUID) (*entities.Account, error) {
	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return nil, err
	}

	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": organizationID})
	if err != nil {
		return nil, err
	}

	if err = s.repo.Assign(ctx, account, organization); err != nil {
		return nil, err
	}

	account.Organizations = append(account.Organizations, organization)

	return account, nil
}

func (s *OrganizationService) Revoke(ctx context.Context, accountID, organizationID uuid.UUID) error {
	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return err
	}

	for _, role := range account.Roles {
		if role.Type == entities.RootRoleTypeName {
			return fmt.Errorf("%v from root user", e.ErrCanNotRemoveOrganization)
		}
	}

	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": organizationID})
	if err != nil {
		return err
	}

	return s.repo.Revoke(ctx, account, organization)
}

func (s *OrganizationService) LoadAccountOrganization(ctx context.Context, account *entities.Account) error {
	organizations, err := s.repo.GetByAccount(ctx, account)
	if err != nil {
		return err
	}

	account.Organizations = organizations

	return nil
}

func (s *OrganizationService) GetIntegratorGames(ctx context.Context, organizationID uuid.UUID) ([]string, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": organizationID})
	if err != nil {
		return nil, err
	}

	if !organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotIntegrator
	}

	return s.gameService.GetIntegratorGameNames(ctx, organization.ID)
}

func (s *OrganizationService) GetByApiKey(ctx context.Context, apiKey string) (*entities.Organization, error) {
	return s.repo.Get(ctx, map[string]interface{}{"api_key": apiKey})
}

func (s *OrganizationService) GetByName(ctx context.Context, name string) (*entities.Organization, error) {
	return s.repo.Get(ctx, map[string]interface{}{"name": name})
}

func (s *OrganizationService) HasAccess(ctx context.Context, organizationID uuid.UUID, gameName string) error {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": organizationID})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	if !organization.IsIntegrator() {
		return e.ErrOrganizationIsNotIntegrator
	}

	games, err := s.gameService.GetIntegratorGameNames(ctx, organizationID)
	if err != nil {
		zap.S().Error(err)

		return err
	}

	for _, game := range games {
		if game == gameName {
			return nil
		}
	}

	return e.ErrDoesNotHavePermission
}

func (s *OrganizationService) GetOrganizationPair(ctx context.Context, providerID, integratorID uuid.UUID) (
	*entities.ProviderIntegratorPair, error) {
	return s.repo.GetOrganizationPair(ctx, providerID, integratorID)
}

func (s *OrganizationService) CreateOrganizationPair(ctx context.Context, providerID, integratorID uuid.UUID) (
	*entities.ProviderIntegratorPair, error) {
	_, err := s.repo.GetOrganizationPair(ctx, providerID, integratorID)
	if err == nil {
		return nil, e.ErrEntityAlreadyExist
	}

	if !errors.Is(err, e.ErrEntityNotFound) {
		return nil, err
	}

	pair := &entities.ProviderIntegratorPair{
		ID:           uuid.New(),
		ProviderID:   providerID,
		IntegratorID: integratorID,
	}

	return s.repo.CreateOrganizationPair(ctx, pair)
}

func (s *OrganizationService) DeleteOrganizationPair(ctx context.Context, providerID, integratorID uuid.UUID) error {
	pair, err := s.repo.GetOrganizationPair(ctx, providerID, integratorID)
	if err != nil {
		return err
	}

	return s.repo.DeleteOrganizationPair(ctx, pair)
}

func (s *OrganizationService) GetIntegratorNames(ctx context.Context, account *entities.Account) ([]string, error) {
	orgs, err := s.repo.GetByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(orgs, func(item *entities.Organization, index int) (string, bool) {
		return item.Name, item.IsIntegrator()
	}), nil
}

func (s *OrganizationService) GetGames(ctx context.Context, integratorID uuid.UUID) ([]*entities.IntegratorGame, error) {
	ig, err := s.repo.GetIntegratorGameList(ctx, integratorID)
	if err != nil {
		return nil, err
	}
	for i := range ig {
		availableWagerSets, err := s.gameService.GetAvailableWagerSetsByIDs(ctx, ig[i].Game)
		if err != nil {
			return nil, err
		}
		ig[i].Game.AvailableWagerSets = availableWagerSets
	}
	return ig, err
}

func (s *OrganizationService) AssignGames(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, gameIDs ...uuid.UUID) ([]*entities.IntegratorGame, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if !organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotIntegrator
	}

	err = s.repo.AssignGames(ctx, integratorID, wagerSetID, gameIDs)
	if err != nil {
		return nil, err
	}

	return s.repo.GetIntegratorGameList(ctx, integratorID)
}

func (s *OrganizationService) UpdateGame(ctx context.Context, integratorID uuid.UUID, gameID uuid.UUID, wagerSetID uuid.UUID, rtp *int64, volatility *string, shortLink bool) (*entities.IntegratorGame, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if !organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotIntegrator
	}

	ig, err := s.repo.UpdateIntegratorGame(ctx, integratorID, gameID, wagerSetID, rtp, volatility, shortLink)
	if err != nil {
		return nil, err
	}

	return ig, nil
}

func (s *OrganizationService) GetIntegratorGameSettings(ctx context.Context, integratorID uuid.UUID, gameName string, currency string) (*entities.IntegratorGame, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if !organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotIntegrator
	}

	games, err := s.gameService.GetIntegratorGames(ctx, organization.ID)
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	var gameID *uuid.UUID
	for _, game := range games {
		if game.Name == gameName {
			gameID = &game.ID
			break
		}
	}

	if gameID == nil {
		return nil, e.ErrGameNotExist
	}

	ig, err := s.repo.GetIntegratorGameSettings(ctx, organization.ID, *gameID, currency)
	if err != nil {
		return nil, err
	}

	return ig, nil
}

func (s *OrganizationService) RevokeGames(ctx context.Context, integratorID uuid.UUID, gameIDs ...uuid.UUID) error {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	if !organization.IsIntegrator() {
		return e.ErrOrganizationIsNotIntegrator
	}

	err = s.repo.RevokeGames(ctx, integratorID, gameIDs)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrganizationService) GetOperatorPair(ctx context.Context, integratorID, operatorID uuid.UUID) (
	*entities.IntegratorOperatorPair, error) {
	return s.repo.GetOperatorPair(ctx, integratorID, operatorID)
}

func (s *OrganizationService) CreateOperatorPair(ctx context.Context, integratorID, operatorID uuid.UUID) (
	*entities.IntegratorOperatorPair, error) {
	_, err := s.repo.GetOperatorPair(ctx, integratorID, operatorID)
	if err == nil {
		return nil, e.ErrEntityAlreadyExist
	}

	if !errors.Is(err, e.ErrEntityNotFound) {
		return nil, err
	}

	pair := &entities.IntegratorOperatorPair{
		ID:           uuid.New(),
		IntegratorID: integratorID,
		OperatorID:   operatorID,
	}

	return s.repo.CreateOperatorPair(ctx, pair)
}

func (s *OrganizationService) DeleteOperatorPair(ctx context.Context, integratorID, operatorID uuid.UUID) error {
	pair, err := s.repo.GetOperatorPair(ctx, integratorID, operatorID)
	if err != nil {
		return err
	}

	return s.repo.DeleteOperatorPair(ctx, pair)
}

func (s *OrganizationService) AssignOperator(ctx context.Context, accountID, operatorID string, integratorID uuid.UUID) (*entities.Account, error) {
	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return nil, err
	}

	operator, err := s.repo.Get(ctx, map[string]interface{}{"id": operatorID})
	if err != nil {
		return nil, err
	}

	if !operator.IsOperator() {
		return nil, e.ErrOrganizationIsNotOperator
	}

	_, err = s.GetOperatorPair(ctx, integratorID, operator.ID)
	if err != nil {
		return nil, err
	}

	if err = s.repo.AssignOperator(ctx, account, operator); err != nil {
		return nil, err
	}

	account.Operators = append(account.Operators, operator)

	return account, nil
}

func (s *OrganizationService) RevokeOperator(ctx context.Context, accountID string, operatorID uuid.UUID, integratorID uuid.UUID) error {
	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return err
	}

	for _, role := range account.Roles {
		if role.Type == entities.RootRoleTypeName {
			return fmt.Errorf("%v from root user", e.ErrCanNotRemoveOrganization)
		}
	}

	operator, err := s.repo.Get(ctx, map[string]interface{}{"id": operatorID})
	if err != nil {
		return err
	}

	if !operator.IsOperator() {
		return e.ErrOrganizationIsNotOperator
	}

	_, err = s.GetOperatorPair(ctx, integratorID, operator.ID)
	if err != nil {
		return err
	}

	return s.repo.RevokeOperator(ctx, account, operator)
}

func (s *OrganizationService) GetOperatorsByIntegrator(ctx context.Context, integratorID uuid.UUID) (operators []*entities.Organization, err error) {
	return s.repo.GetOperatorsByIntegrator(ctx, integratorID)
}

func (s *OrganizationService) GetProvidersByIntegrator(ctx context.Context, integratorID uuid.UUID) (providers []*entities.Organization, err error) {
	return s.repo.GetProvidersByIntegrator(ctx, integratorID)
}

func (s *OrganizationService) GetOrganizationPairsByIntegrator(ctx context.Context, integratorID uuid.UUID) (
	pairIDs []uuid.UUID, err error) {
	providerIntegratorPairs, err := s.repo.GetOrganizationPairsByIntegrator(ctx, integratorID)
	if err != nil {
		return
	}
	for _, pair := range providerIntegratorPairs {
		pairIDs = append(pairIDs, pair.ID)
	}

	return
}

func (s *OrganizationService) GetGamesWagerSets(ctx context.Context, integratorID uuid.UUID) ([]*entities.IntegratorGameWagerSet, error) {
	igws, err := s.repo.GetIntegratorGameWagerSetList(ctx, integratorID)
	if err != nil {
		return nil, err
	}

	return igws, err
}

func (s *OrganizationService) CreateIntegratorGameWagerSet(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, currency string, gameID uuid.UUID) ([]*entities.IntegratorGameWagerSet, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if !organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotIntegrator
	}

	err = s.repo.CreateIntegratorGameWagerSet(ctx, integratorID, wagerSetID, currency, gameID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetIntegratorGameWagerSetList(ctx, integratorID)
}

func (s *OrganizationService) UpdateGameWagerSet(ctx context.Context,
	integratorID uuid.UUID, gameID uuid.UUID,
	wagerSetID uuid.UUID, currency string,
	newCurrency string) (*entities.IntegratorGameWagerSet, error) {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return nil, err
	}

	if !organization.IsIntegrator() {
		return nil, e.ErrOrganizationIsNotIntegrator
	}

	igws, err := s.repo.UpdateIntegratorGameWagerSet(ctx, integratorID, gameID, wagerSetID, currency, newCurrency)
	if err != nil {
		return nil, err
	}

	return igws, nil
}

func (s *OrganizationService) DeleteGameWagerSet(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, currency string, gameID uuid.UUID) error {
	organization, err := s.repo.Get(ctx, map[string]interface{}{"id": integratorID})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	if !organization.IsIntegrator() {
		return e.ErrOrganizationIsNotIntegrator
	}

	err = s.repo.DeleteGameWagerSet(ctx, integratorID, wagerSetID, currency, gameID)
	if err != nil {
		return err
	}

	return nil
}
