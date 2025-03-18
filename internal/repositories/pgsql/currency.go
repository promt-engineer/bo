package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type currencyRepository struct {
	conn *gorm.DB
}

func NewCurrencyRepository(conn *gorm.DB) *currencyRepository {
	return &currencyRepository{conn: conn}
}

func (r *currencyRepository) All(ctx context.Context) ([]*entities.CurrencyMultiplier, error) {
	res := []*entities.CurrencyMultiplier{}

	return res, r.conn.WithContext(ctx).
		Joins("ProviderIntegratorPair").
		Preload("ProviderIntegratorPair.Provider").
		Preload("ProviderIntegratorPair.Integrator").
		Find(&res).Error
}

func (r *currencyRepository) UniqueCurrencyNames(ctx context.Context) ([]string, error) {
	var uniqueNames []string

	if err := r.conn.WithContext(ctx).
		Model(&entities.CurrencyMultiplier{}).
		Distinct("title").
		Pluck("title", &uniqueNames).
		Error; err != nil {
		return nil, err
	}

	return uniqueNames, nil
}

func (r *currencyRepository) Search(ctx context.Context, filter map[string]interface{}) (cm []*entities.CurrencyMultiplier, err error) {
	query := r.conn.WithContext(ctx).
		Joins("ProviderIntegratorPair").
		Preload("ProviderIntegratorPair.Provider").
		Preload("ProviderIntegratorPair.Integrator").
		Model(&entities.CurrencyMultiplier{})

	if len(filter) > 0 {
		//query = query.Where(filter)
		for key, value := range filter {
			switch v := value.(type) {
			case []uuid.UUID:
				query = query.Where(key+" IN (?)", v)
			default:
				query = query.Where(key+" = ?", v)
			}
		}
	}

	if err = query.Find(&cm).Error; err != nil {
		return
	}

	return
}

func (r *currencyRepository) Get(ctx context.Context, params map[string]interface{}) (cm *entities.CurrencyMultiplier, err error) {
	if err = r.conn.WithContext(ctx).Where(params).
		Preload("ProviderIntegratorPair").
		Preload("ProviderIntegratorPair.Provider").
		Preload("ProviderIntegratorPair.Integrator").
		First(&cm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	return
}

func (r *currencyRepository) CreateCurrencyMultiplier(ctx context.Context, cm *entities.CurrencyMultiplier) (*entities.CurrencyMultiplier, error) {
	if err := r.conn.WithContext(ctx).Create(cm).Error; err != nil {
		return nil, err
	}

	return r.Get(ctx, map[string]interface{}{"organization_pair_id": cm.OrganizationPairID, "title": cm.Title})
}

func (r *currencyRepository) UpdateCurrencyMultiplier(ctx context.Context, cm *entities.CurrencyMultiplier) (*entities.CurrencyMultiplier, error) {
	if err := r.conn.WithContext(ctx).
		Where("organization_pair_id = ? and title = ?", cm.OrganizationPairID, cm.Title).
		Updates(cm).Error; err != nil {
		return nil, err
	}

	return r.Get(ctx, map[string]interface{}{"organization_pair_id": cm.OrganizationPairID, "title": cm.Title})
}

func (r *currencyRepository) DeleteCurrencyMultiplier(ctx context.Context, cm *entities.CurrencyMultiplier) error {
	return r.conn.WithContext(ctx).
		Where("organization_pair_id = ? and title = ?", cm.OrganizationPairID, cm.Title).Delete(&cm).Error
}

func (r *currencyRepository) CurrencyGetAll(ctx context.Context, filters map[string]interface{}) (currencies []*entities.Currency, err error) {
	err = r.conn.WithContext(ctx).Where(filters).Find(&currencies).Error

	return
}

func (r *currencyRepository) CurrencyGet(ctx context.Context, alias string) (currency *entities.Currency, err error) {
	err = r.conn.WithContext(ctx).Where("alias = ?", alias).First(&currency).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}
		return nil, err
	}
	return currency, nil
}

func (r *currencyRepository) CreateCurrency(ctx context.Context, currency *entities.Currency) (*entities.Currency, error) {
	if err := r.conn.WithContext(ctx).Create(currency).Error; err != nil {
		return nil, err
	}

	return r.CurrencyGet(ctx, currency.Title)
}

func (r *currencyRepository) DeleteCurrency(ctx context.Context, currency *entities.Currency) error {
	return r.conn.WithContext(ctx).
		Where("alias = ?", currency.Alias).Delete(&currency).Error
}

//func (r *currencyRepository) GetCurrencyMultiplier(ctx context.Context, filters map[string]interface{}) (cm []*entities.CurrencyMultiplier, err error) {
//	query := r.conn.WithContext(ctx)
//
//	for key, value := range filters {
//		switch v := value.(type) {
//		case []uuid.UUID:
//			query = query.Where(key+" IN (?)", v)
//		default:
//			query = query.Where(key+" = ?", v)
//		}
//	}
//
//	err = query.Find(&cm).Error
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, e.ErrEntityNotFound
//		}
//		return nil, err
//	}
//	return cm, nil
//}
