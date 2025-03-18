package requests

import (
	"github.com/google/uuid"
	"time"
)

type CurrencyMultiplier struct {
	CurrencyMultiplierIdentify
	Multiplier int64  `json:"multiplier" form:"multiplier"`
	Synonym    string `json:"synonym" form:"synonym" validate:"min=1"`
}

type CurrencyConfig struct {
	ProviderID   uuid.UUID `json:"provider_id" form:"provider_id"`
	IntegratorID uuid.UUID `json:"integrator_id" form:"integrator_id"`
	DefaultWager int64     `json:"default_wager" form:"default_wager"`
	WagerLevels  []int64   `json:"wager_levels" form:"wager_levels"`
}

type CurrencySearchRequest struct {
	OrganizationPairID uuid.UUID `json:"organization_pair_id" form:"organization_pair_id"`
}

type CurrencyOrganizationPair struct {
	ProviderID   uuid.UUID `json:"provider_id" form:"provider_id"`
	IntegratorID uuid.UUID `json:"integrator_id" form:"integrator_id"`
}

type CurrencyMultiplierIdentify struct {
	OrganizationPairID uuid.UUID `json:"organization_pair_id" form:"organization_pair_id"`
	Title              string    `json:"title" form:"title" validate:"min=1"`
}

type Currency struct {
	Title        string  `json:"title"  form:"title" validate:"min=1"`
	Alias        string  `json:"alias" form:"alias" validate:"min=1"`
	Type         string  `json:"type" form:"type" validate:"min=1"`
	Rate         float64 `json:"additional_alias_rate" form:"additional_alias_rate" validate:"required,gte=1"`
	BaseCurrency string  `json:"base_currency,omitempty" form:"base_currency"`
}

type CurrencyExchangeRequest struct {
	From string  `json:"from" validate:"required"`
	Rate float64 `json:"rate" validate:"required"`
}

type DeleteCurrencyExchangeRequest struct {
	CurrencyExchangeRequest
	CreatedAt time.Time `json:"created_at"`
}
