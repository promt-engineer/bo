package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type organizationRepository struct {
	conn *gorm.DB
}

func NewOrganizationRepository(conn *gorm.DB) *organizationRepository {
	return &organizationRepository{
		conn: conn,
	}
}

func (r *organizationRepository) All(ctx context.Context) (organizations []*entities.Organization, err error) {
	err = r.conn.WithContext(ctx).Find(&organizations).Error

	return
}

func (r *organizationRepository) GetIntegratorsByProvider(ctx context.Context, providerID uuid.UUID) (organizations []*entities.Organization, err error) {
	err = r.conn.WithContext(ctx).
		Joins("join integrator_providers on integrator_providers.integrator_id = organizations.id").
		Where("integrator_providers.provider_id = ?", providerID).
		Find(&organizations).Error

	return
}

func (r *organizationRepository) Create(ctx context.Context, organization *entities.Organization) (*entities.Organization, error) {
	if err := r.conn.WithContext(ctx).Where("name = ? and type = ?", organization.Name, organization.Type).First(&entities.Organization{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = r.conn.WithContext(ctx).Create(&organization).Error; err != nil {
				return nil, err
			}

			return organization, nil
		}

		return nil, err
	}

	return nil, e.ErrOrganizationNameMustBeUnique
}

func (r *organizationRepository) Update(ctx context.Context, organization *entities.Organization) (*entities.Organization, error) {
	if err := r.conn.WithContext(ctx).Updates(&organization).Error; err != nil {
		return nil, err
	}

	return r.Get(ctx, map[string]interface{}{"id": organization.ID})
}

func (r *organizationRepository) Delete(ctx context.Context, accountID uuid.UUID, organization *entities.Organization) error {
	var res []entities.AccountOrganization
	err := r.conn.WithContext(ctx).Where("organization_id = ?", organization.ID).Find(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if len(res) > 1 {
		return e.ErrOrganizationInUse
	}

	if len(res) == 1 {
		if res[0].AccountID != accountID {
			return e.ErrOrganizationInUse
		}
	}

	return r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(res) == 1 {
			if err := r.conn.WithContext(ctx).Where("organization_id = ?", organization.ID).Delete(&entities.AccountOrganization{}).Error; err != nil {
				return err
			}
		}

		if err := r.conn.WithContext(ctx).Where("organization_id = ?", organization.ID).Delete(&entities.Role{}).Error; err != nil {
			return err
		}

		if err := r.conn.WithContext(ctx).Where("id = ?", organization.ID).Delete(&organization).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *organizationRepository) Get(ctx context.Context, params map[string]interface{}) (organization *entities.Organization, err error) {
	if err = r.conn.WithContext(ctx).Where(params).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	return
}

func (r *organizationRepository) Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) (organization []*entities.Organization, total int64, err error) {
	query := r.conn.WithContext(ctx).Model(&entities.Organization{}).Where(filters)

	if err = query.Count(&total).Error; err != nil {
		return
	}

	if err = query.Order(order).Limit(limit).Offset(offset).Find(&organization).Error; err != nil {
		return
	}
	return
}

func (r *organizationRepository) Assign(ctx context.Context, account *entities.Account, organization *entities.Organization) error {
	var ao *entities.AccountOrganization

	err := r.conn.WithContext(ctx).Where("account_id = ? and organization_id = ?", account.ID, organization.ID).First(&ao).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return e.ErrOrganizationAlreadyAssigned
	}

	return r.conn.WithContext(ctx).Create(&entities.AccountOrganization{AccountID: account.ID, OrganizationID: organization.ID}).Error
}

func (r *organizationRepository) Revoke(ctx context.Context, account *entities.Account, organization *entities.Organization) error {
	return r.conn.WithContext(ctx).Where("account_id = ? and organization_id = ?", account.ID, organization.ID).Delete(&entities.AccountOrganization{}).Error
}

func (r *organizationRepository) GetByAccount(ctx context.Context, account *entities.Account) ([]*entities.Organization, error) {
	var organizations []*entities.Organization

	if err := r.conn.WithContext(ctx).
		Joins("join account_organizations on account_organizations.organization_id = organizations.id and account_organizations.account_id = ?", account.ID).
		Find(&organizations).Error; err != nil {
		return nil, err
	}

	return organizations, nil
}

func (r *organizationRepository) IntegratorGameExists(ctx context.Context, organization *entities.Organization, game *entities.Game) bool {
	var ig *entities.IntegratorGame

	err := r.conn.WithContext(ctx).Where("organization_id = ? and game_id = ?", organization.ID, game.ID).Find(&ig).Error
	if err != nil {
		zap.S().Error(err)

		return false
	}

	return ig != nil
}

func (r *organizationRepository) GetOrganizationPair(ctx context.Context, providerID, integratorID uuid.UUID) (
	pair *entities.ProviderIntegratorPair, err error) {
	if err = r.conn.WithContext(ctx).
		Where("provider_id = ? and integrator_id = ?", providerID, integratorID).
		First(&pair).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pair, e.ErrEntityNotFound
		}

		return pair, err
	}

	return
}

func (r *organizationRepository) CreateOrganizationPair(ctx context.Context, pair *entities.ProviderIntegratorPair) (
	p *entities.ProviderIntegratorPair, err error) {
	err = r.conn.WithContext(ctx).Create(pair).Error
	if err != nil {
		return p, err
	}

	return r.GetOrganizationPair(ctx, pair.ProviderID, pair.IntegratorID)
}

func (r *organizationRepository) DeleteOrganizationPair(ctx context.Context, pair *entities.ProviderIntegratorPair) error {
	return r.conn.WithContext(ctx).Where("id = ?", pair.ID).Delete(&pair).Error
}

func (r *organizationRepository) GetIntegratorGameList(ctx context.Context, integratorID uuid.UUID) (ig []*entities.IntegratorGame, err error) {
	err = r.conn.WithContext(ctx).
		Preload("Organization").
		Preload("Game").
		Preload("WagerSet").
		Where("organization_id = ?", integratorID).
		Find(&ig).Error

	return ig, err
}

func (r *organizationRepository) AssignGames(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, gameIDs []uuid.UUID) error {
	for _, gameID := range gameIDs {
		integratorGame := entities.IntegratorGame{
			OrganizationID: integratorID,
			Organization:   nil,
			GameID:         gameID,
			Game:           nil,
			WagerSetID:     wagerSetID,
			WagerSet:       nil,
		}

		err := r.conn.WithContext(ctx).Create(&integratorGame).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *organizationRepository) UpdateIntegratorGame(ctx context.Context, integratorID uuid.UUID, gameID uuid.UUID, wagerSetID uuid.UUID, rtp *int64, volatility *string, shortLink bool) (ig *entities.IntegratorGame, err error) {
	err = r.conn.WithContext(ctx).Where("organization_id = ? AND game_id = ?", integratorID, gameID).First(&ig).Error
	if err != nil {
		return nil, err
	}

	ig.WagerSetID = wagerSetID
	ig.RTP = rtp
	ig.Volatility = volatility
	ig.ShortLink = shortLink

	err = r.conn.WithContext(ctx).Where("organization_id = ? AND game_id = ?", integratorID, gameID).Updates(&ig).Error
	if err != nil {
		return nil, err
	}

	return ig, nil
}

func (r *organizationRepository) RevokeGames(ctx context.Context, integratorID uuid.UUID, gameIDs []uuid.UUID) error {
	for _, gameID := range gameIDs {
		err := r.conn.WithContext(ctx).Where("organization_id = ? AND game_id = ?", integratorID, gameID).Delete(&entities.IntegratorGame{}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *organizationRepository) GetIntegratorGameSettings(ctx context.Context, integratorID, gameID uuid.UUID, currency string) (ig *entities.IntegratorGame, err error) {
	query := `
		SELECT 
			ig.organization_id, 
			ig.game_id, 
			COALESCE(gmws.wager_set_id, ig.wager_set_id) AS wager_set_id, 
			ig.rtp, 
			ig.volatility, 
			ig.short_link
		FROM public.integrator_games ig
		LEFT JOIN public.integrator_game_wager_sets gmws 
			ON ig.organization_id = gmws.organization_id 
			AND ig.game_id = gmws.game_id 
			AND gmws.currency = ?
		WHERE ig.organization_id = ? AND ig.game_id = ?
	`

	err = r.conn.WithContext(ctx).
		Raw(query, currency, integratorID, gameID).
		Scan(&ig).Error
	if err != nil {
		return nil, err
	}

	if err = r.conn.WithContext(ctx).Model(&ig).Association("WagerSet").Find(&ig.WagerSet); err != nil {
		return nil, err
	}

	return ig, nil
}

func (r *organizationRepository) GetOperatorPair(ctx context.Context, integratorID, operatorID uuid.UUID) (
	pair *entities.IntegratorOperatorPair, err error) {
	if err = r.conn.WithContext(ctx).
		Where("integrator_id = ? and operator_id = ?", integratorID, operatorID).
		First(&pair).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pair, e.ErrEntityNotFound
		}

		return pair, err
	}

	return
}

func (r *organizationRepository) CreateOperatorPair(ctx context.Context, pair *entities.IntegratorOperatorPair) (
	p *entities.IntegratorOperatorPair, err error) {
	err = r.conn.WithContext(ctx).Create(pair).Error
	if err != nil {
		return p, err
	}

	return r.GetOperatorPair(ctx, pair.IntegratorID, pair.OperatorID)
}

func (r *organizationRepository) DeleteOperatorPair(ctx context.Context, pair *entities.IntegratorOperatorPair) error {
	return r.conn.WithContext(ctx).Where("id = ?", pair.ID).Delete(&pair).Error
}

func (r *organizationRepository) AssignOperator(ctx context.Context, account *entities.Account, operator *entities.Organization) error {
	var ao *entities.AccountOperator

	err := r.conn.WithContext(ctx).Where("account_id = ? and operator_id = ?", account.ID, operator.ID).First(&ao).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return e.ErrOperatorAlreadyAssigned
	}

	return r.conn.WithContext(ctx).Create(&entities.AccountOperator{AccountID: account.ID, OperatorID: operator.ID}).Error
}

func (r *organizationRepository) RevokeOperator(ctx context.Context, account *entities.Account, operator *entities.Organization) error {
	return r.conn.WithContext(ctx).Where("account_id = ? and operator_id = ?", account.ID, operator.ID).Delete(&entities.AccountOperator{}).Error
}

func (r *organizationRepository) GetOperatorsByIntegrator(ctx context.Context, integratorID uuid.UUID) (operators []*entities.Organization, err error) {
	err = r.conn.WithContext(ctx).
		Joins("join operator_integrators on operator_integrators.operator_id = organizations.id").
		Where("operator_integrators.integrator_id = ?", integratorID).
		Find(&operators).Error

	return
}

func (r *organizationRepository) GetProvidersByIntegrator(ctx context.Context, integratorID uuid.UUID) (providers []*entities.Organization, err error) {
	err = r.conn.WithContext(ctx).
		Joins("join integrator_providers on integrator_providers.provider_id = organizations.id").
		Where("integrator_providers.integrator_id = ?", integratorID).
		Find(&providers).Error

	return
}

func (r *organizationRepository) GetOrganizationPairsByIntegrator(ctx context.Context, integratorID uuid.UUID) (pairs []*entities.ProviderIntegratorPair, err error) {
	if err = r.conn.WithContext(ctx).
		Where("integrator_id = ?", integratorID).
		Find(&pairs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pairs, e.ErrEntityNotFound
		}

		return pairs, err
	}

	return
}

func (r *organizationRepository) GetIntegratorGameWagerSetList(ctx context.Context, integratorID uuid.UUID) (igws []*entities.IntegratorGameWagerSet, err error) {
	err = r.conn.WithContext(ctx).
		Preload("Organization").
		Preload("Game").
		Preload("WagerSet").
		Where("organization_id = ?", integratorID).
		Find(&igws).Error

	return igws, err
}

func (r *organizationRepository) CreateIntegratorGameWagerSet(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, currency string, gameID uuid.UUID) error {
	integratorGameWagerSet := entities.IntegratorGameWagerSet{
		OrganizationID: integratorID,
		GameID:         gameID,
		WagerSetID:     wagerSetID,
		Currency:       currency,
	}

	err := r.conn.WithContext(ctx).Create(&integratorGameWagerSet).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *organizationRepository) UpdateIntegratorGameWagerSet(ctx context.Context, integratorID uuid.UUID, gameID uuid.UUID, wagerSetID uuid.UUID, currency string, newCurrency string) (igws *entities.IntegratorGameWagerSet, err error) {
	err = r.conn.WithContext(ctx).
		Preload("Organization").
		Preload("Game").
		Preload("WagerSet").
		Where("organization_id = ? AND game_id = ? AND wager_set_id = ? AND currency = ?", integratorID, gameID, wagerSetID, currency).First(&igws).Error
	if err != nil {
		return nil, err
	}

	if newCurrency != "" {
		igws.Currency = newCurrency
	}

	err = r.conn.WithContext(ctx).
		Where("organization_id = ? AND game_id = ? AND wager_set_id = ? AND currency = ?", integratorID, gameID, wagerSetID, currency).Updates(&igws).Error
	if err != nil {
		return nil, err
	}

	return igws, nil
}

func (r *organizationRepository) DeleteGameWagerSet(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, currency string, gameID uuid.UUID) error {
	err := r.conn.WithContext(ctx).Where("organization_id = ? AND game_id = ? AND wager_set_id = ? AND currency = ?", integratorID, gameID, wagerSetID, currency).Delete(&entities.IntegratorGameWagerSet{}).Error

	if err != nil {
		return err
	}

	return nil
}
