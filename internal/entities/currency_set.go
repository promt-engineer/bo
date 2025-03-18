package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"strings"
	"time"
)

type CurrencySet struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ID             uuid.UUID      `json:"id"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	Name           string         `json:"name"`
	Currencies     pq.StringArray `json:"currencies" gorm:"type:varchar[]" swaggertype:"array,string"`

	IsActive bool `json:"is_active"`
}

func (*CurrencySet) TableName() string {
	return "currency_sets"
}

func (cs *CurrencySet) SetCurrencies(allCurrencies, currencies []string) {
	currencies = lo.Map(currencies, func(item string, index int) string {
		return strings.ToLower(item)
	})

	currencies = lo.Uniq(currencies)

	currencies = lo.Filter(currencies, func(item string, index int) bool {
		return lo.Contains(allCurrencies, item)
	})

	cs.Currencies = currencies
}
