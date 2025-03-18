package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type gameRepository struct {
	conn *gorm.DB
}

func NewGameRepository(conn *gorm.DB) *gameRepository {
	return &gameRepository{conn: conn}
}

func (r *gameRepository) GetBy(ctx context.Context, condition map[string]interface{}) (game *entities.Game, err error) {
	if err = r.conn.WithContext(ctx).Preload("Organization").Where(condition).First(&game).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	game.AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)

	return
}

func (r *gameRepository) Create(ctx context.Context, game *entities.Game) (*entities.Game, error) {
	err := r.conn.WithContext(ctx).Create(game).Error
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameRepository) All(ctx context.Context, organizationID *uuid.UUID, condition map[string]interface{}) (games []*entities.Game, err error) {
	conn := r.conn.WithContext(ctx)
	conn = r.withOwner(conn, organizationID).Preload("WagerSet")

	if err = conn.Where(condition).Find(&games).Error; err != nil {
		return nil, err
	}

	for i, game := range games {
		games[i].AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)
		if err != nil {
			return nil, err
		}
	}

	return games, nil
}

func (r *gameRepository) Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) (games []*entities.Game, total int64, err error) {
	query := r.conn.WithContext(ctx).Model(&entities.Game{}).Where(filters)

	if err = query.Count(&total).Error; err != nil {
		return
	}

	if err = query.Order(order).Limit(limit).Offset(offset).Preload("WagerSet").Find(&games).Error; err != nil {
		return
	}
	for i, game := range games {
		games[i].AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)
		if err != nil {
			return
		}
	}
	return
}

func (r *gameRepository) GetDictionaries(ctx context.Context, organizationID *uuid.UUID, dictType string) (dict []string, err error) {
	conn := r.conn.WithContext(ctx)
	conn = r.withOwner(conn, organizationID).Preload("WagerSet")

	query := fmt.Sprintf(`SELECT DISTINCT unnest(%s) FROM games`, dictType)

	rows, err := conn.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		dict = append(dict, value)
	}

	return dict, nil
}

func (r *gameRepository) GetAllByFilter(ctx context.Context, condition map[string]interface{}) (games []*entities.Game, err error) {
	query := r.conn.WithContext(ctx).Model(&entities.Game{})

	if len(condition) > 0 {
		query = query.Where(condition)
	}

	if err = query.Find(&games).Error; err != nil {
		return
	}

	for i, game := range games {
		games[i].AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (r *gameRepository) GetOrganizationGameList(ctx context.Context, organizationID uuid.UUID) (games []*entities.Game, err error) {
	err = r.conn.WithContext(ctx).
		Select("distinct (games.id) as _, games.*").
		Joins(`inner join integrator_providers as ip 
						on games.organization_id = ip.provider_id`).
		Where("ip.integrator_id = ? or ip.provider_id = ?", organizationID, organizationID).
		Find(&games).
		Error

	for i, game := range games {
		games[i].AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (r *gameRepository) GetIntegratorGameList(ctx context.Context, organizationID uuid.UUID) (games []*entities.Game, err error) {
	err = r.conn.WithContext(ctx).
		Select("DISTINCT games.id AS _, games.*, COALESCE(NULLIF(ig.wager_set_id, '00000000-0000-0000-0000-000000000000'), games.wager_set_id) AS wager_set_id").
		Joins(`INNER JOIN integrator_games AS ig ON games.id = ig.game_id`).
		Joins(`LEFT JOIN wager_sets AS ws ON ws.id = COALESCE(NULLIF(ig.wager_set_id, '00000000-0000-0000-0000-000000000000'), games.wager_set_id)`).
		Where("ig.organization_id = ?", organizationID).
		Preload("WagerSet").
		Find(&games).
		Error
	if err != nil {
		return nil, err
	}

	for i, game := range games {
		games[i].AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)
		if err != nil {
			return nil, err
		}
	}

	return games, nil
}

func (r *gameRepository) Update(ctx context.Context, gameID uuid.UUID, condition map[string]interface{}) (*entities.Game, error) {
	if err := r.conn.WithContext(ctx).Model(&entities.Game{}).Where("id = ?", gameID).Updates(condition).Error; err != nil {
		return nil, err
	}

	return r.GetBy(ctx, map[string]interface{}{"id": gameID})
}

func (r *gameRepository) Delete(ctx context.Context, game *entities.Game) error {
	return r.conn.WithContext(ctx).Where("id = ?", game.ID).Delete(&game).Error
}

func (r *gameRepository) withOwner(conn *gorm.DB, organizationID *uuid.UUID) *gorm.DB {
	conn = conn.Joins("Organization")

	if organizationID != nil {
		conn = conn.Where("organization_id = ?", organizationID)
	}

	return conn
}

func (r *gameRepository) AddValueToTheDictionary(ctx context.Context, organizationID *uuid.UUID, dictType, value string) (string, error) {
	conn := r.conn.WithContext(ctx)

	query := fmt.Sprintf("UPDATE games SET %s = array_append(%s, $1) WHERE organization_id = $2", dictType, dictType)

	tx := conn.Exec(query, value, organizationID)
	if tx.Error != nil {
		return "", tx.Error
	}

	return value, nil
}

func (r *gameRepository) RemoveValueFromDictionary(ctx context.Context, organizationID *uuid.UUID, dictType, value string) error {
	conn := r.conn.WithContext(ctx)

	query := fmt.Sprintf("UPDATE games SET %s = array_remove(%s, $1) WHERE organization_id = $2", dictType, dictType)

	tx := conn.Exec(query, value, organizationID)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *gameRepository) GetIntegratorGame(ctx context.Context, organizationID uuid.UUID, gameName string) (game *entities.Game, err error) {
	err = r.conn.WithContext(ctx).
		Joins("INNER JOIN integrator_games AS ig ON games.id = ig.game_id").
		Where("ig.organization_id = ? AND games.name = ?", organizationID, gameName).
		First(&game).Error
	if err != nil {
		return nil, err
	}

	game.AvailableWagerSets, err = r.GetAvailableWagerSetsByIDs(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameRepository) GetAvailableWagerSetsByIDs(ctx context.Context, game *entities.Game) (wagerSets []entities.WagerSet, err error) {
	if len(game.AvailableWagerSetsID) > 0 {
		err = r.conn.WithContext(ctx).Raw(`
				SELECT * FROM wager_sets WHERE id = ANY(?)
			`, pq.Array(game.AvailableWagerSetsID)).Scan(&wagerSets).Error
		if err != nil {
			return
		}
	}

	return
}

func (r *gameRepository) GetWagerSetByID(ctx context.Context, id uuid.UUID) (*entities.WagerSet, error) {
	var wagerSet entities.WagerSet
	if err := r.conn.Where("id = ?", id).First(&wagerSet).Error; err != nil {
		return nil, err
	}
	return &wagerSet, nil
}
