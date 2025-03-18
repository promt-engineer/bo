package repositories

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
)

type GameRepository interface {
	GetBy(ctx context.Context, condition map[string]interface{}) (*entities.Game, error)
	Create(ctx context.Context, game *entities.Game) (*entities.Game, error)
	All(ctx context.Context, organizationID *uuid.UUID, condition map[string]interface{}) (games []*entities.Game, err error)
	GetAllByFilter(ctx context.Context, condition map[string]interface{}) (games []*entities.Game, err error)
	GetOrganizationGameList(ctx context.Context, organizationID uuid.UUID) ([]*entities.Game, error)
	GetIntegratorGameList(ctx context.Context, organizationID uuid.UUID) ([]*entities.Game, error)
	Update(ctx context.Context, gameID uuid.UUID, condition map[string]interface{}) (*entities.Game, error)
	Delete(ctx context.Context, game *entities.Game) error
	GetDictionaries(ctx context.Context, organizationID *uuid.UUID, dictType string) (dict []string, err error)
	Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) (games []*entities.Game, total int64, err error)
	AddValueToTheDictionary(ctx context.Context, organizationID *uuid.UUID, dictType, locale string) (string, error)
	RemoveValueFromDictionary(ctx context.Context, organizationID *uuid.UUID, dictType, value string) error
	GetIntegratorGame(ctx context.Context, organizationID uuid.UUID, gameName string) (game *entities.Game, err error)
	GetAvailableWagerSetsByIDs(ctx context.Context, game *entities.Game) (wagerSets []entities.WagerSet, err error)
	GetWagerSetByID(ctx context.Context, id uuid.UUID) (*entities.WagerSet, error)
}
