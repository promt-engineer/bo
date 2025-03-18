package repositories

import (
	"backoffice/internal/entities"
	"context"
)

type CurrencyRepository interface {
	All(ctx context.Context) ([]*entities.CurrencyMultiplier, error)
	UniqueCurrencyNames(ctx context.Context) ([]string, error)
	Search(ctx context.Context, filter map[string]interface{}) (cm []*entities.CurrencyMultiplier, err error)
	CreateCurrencyMultiplier(ctx context.Context, cm *entities.CurrencyMultiplier) (*entities.CurrencyMultiplier, error)
	Get(ctx context.Context, params map[string]interface{}) (account *entities.CurrencyMultiplier, err error)
	UpdateCurrencyMultiplier(ctx context.Context, cm *entities.CurrencyMultiplier) (*entities.CurrencyMultiplier, error)
	DeleteCurrencyMultiplier(ctx context.Context, cm *entities.CurrencyMultiplier) error
	CurrencyGetAll(ctx context.Context, filters map[string]interface{}) (currencies []*entities.Currency, err error)
	CurrencyGet(ctx context.Context, alias string) (currency *entities.Currency, err error)
	CreateCurrency(ctx context.Context, currency *entities.Currency) (*entities.Currency, error)
	DeleteCurrency(ctx context.Context, currency *entities.Currency) error
	//GetCurrencyMultiplier(ctx context.Context, filters map[string]interface{}) (cm []*entities.CurrencyMultiplier, err error)
}
