package repositories

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
)

type OrganizationRepository interface {
	All(ctx context.Context) (organizations []*entities.Organization, err error)
	Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) (organization []*entities.Organization, total int64, err error)
	Create(ctx context.Context, organization *entities.Organization) (*entities.Organization, error)
	Update(ctx context.Context, organization *entities.Organization) (*entities.Organization, error)
	Delete(ctx context.Context, accountID uuid.UUID, organization *entities.Organization) error
	Get(ctx context.Context, params map[string]interface{}) (organization *entities.Organization, err error)
	Assign(ctx context.Context, account *entities.Account, organization *entities.Organization) error
	Revoke(ctx context.Context, account *entities.Account, organization *entities.Organization) error
	GetByAccount(ctx context.Context, account *entities.Account) ([]*entities.Organization, error)
	IntegratorGameExists(ctx context.Context, organization *entities.Organization, game *entities.Game) bool
	GetIntegratorsByProvider(ctx context.Context, providerID uuid.UUID) (organizations []*entities.Organization, err error)
	GetOrganizationPair(ctx context.Context, providerID, integratorID uuid.UUID) (*entities.ProviderIntegratorPair, error)
	DeleteOrganizationPair(ctx context.Context, pair *entities.ProviderIntegratorPair) error
	CreateOrganizationPair(ctx context.Context, pair *entities.ProviderIntegratorPair) (p *entities.ProviderIntegratorPair, err error)
	AssignGames(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, gameIDs []uuid.UUID) error
	RevokeGames(ctx context.Context, integratorID uuid.UUID, gameIDs []uuid.UUID) error
	GetIntegratorGameList(ctx context.Context, integratorID uuid.UUID) (ig []*entities.IntegratorGame, err error)
	UpdateIntegratorGame(ctx context.Context, integratorID, gameID, wagerSetID uuid.UUID, rtp *int64, volatility *string, shortLink bool) (ig *entities.IntegratorGame, err error)
	GetIntegratorGameSettings(ctx context.Context, integratorID, gameID uuid.UUID, currency string) (ig *entities.IntegratorGame, err error)
	GetOperatorPair(ctx context.Context, integratorID, operatorID uuid.UUID) (*entities.IntegratorOperatorPair, error)
	DeleteOperatorPair(ctx context.Context, pair *entities.IntegratorOperatorPair) error
	CreateOperatorPair(ctx context.Context, pair *entities.IntegratorOperatorPair) (p *entities.IntegratorOperatorPair, err error)
	AssignOperator(ctx context.Context, account *entities.Account, operator *entities.Organization) error
	RevokeOperator(ctx context.Context, account *entities.Account, operator *entities.Organization) error
	GetOperatorsByIntegrator(ctx context.Context, integratorID uuid.UUID) (operators []*entities.Organization, err error)
	GetProvidersByIntegrator(ctx context.Context, integratorID uuid.UUID) (providers []*entities.Organization, err error)
	GetOrganizationPairsByIntegrator(ctx context.Context, integratorID uuid.UUID) ([]*entities.ProviderIntegratorPair, error)
	GetIntegratorGameWagerSetList(ctx context.Context, integratorID uuid.UUID) (igws []*entities.IntegratorGameWagerSet, err error)
	CreateIntegratorGameWagerSet(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, currency string, gameID uuid.UUID) error
	UpdateIntegratorGameWagerSet(ctx context.Context, integratorID, gameID, wagerSetID uuid.UUID, currency string, newCurrency string) (igws *entities.IntegratorGameWagerSet, err error)
	DeleteGameWagerSet(ctx context.Context, integratorID uuid.UUID, wagerSetID uuid.UUID, currency string, gameID uuid.UUID) error
}
