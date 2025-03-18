package entities

import (
	"backoffice/pkg/exchange"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"time"
)

//type CurrencyConfig struct {
//	CreatedAt time.Time      `json:"-"`
//	UpdatedAt time.Time      `json:"-"`
//	DeletedAt gorm.DeletedAt `json:"-"`
//
//	OrganizationPairID     uuid.UUID               `json:"organization_pair_id"`
//	ProviderIntegratorPair *ProviderIntegratorPair `json:"provider_integrator_pair" gorm:"foreignKey:organization_pair_id"`
//
//	DefaultWager int64         `json:"default_wager"`
//	WagerLevels  pq.Int64Array `json:"wager_levels" gorm:"type:int[]" swaggerignore:"true"`
//
//	MultipliersRaw []CurrencyMultiplier `json:"-" gorm:"foreignKey:OrganizationPairID;references:OrganizationPairID"`
//
//	// computed
//	Multipliers map[string]int64 `json:"multipliers" gorm:"-"`
//}

//	func (cc *CurrencyConfig) Currencies() []string {
//		return lo.Keys(cc.Multipliers)
//	}
//
//	func (cc *CurrencyConfig) Compute() {
//		if cc.Multipliers == nil {
//			cc.Multipliers = map[string]int64{}
//		}
//
//		for _, cm := range cc.MultipliersRaw {
//			cc.Multipliers[cm.Title] = cm.Multiplier
//		}
//	}
//
//	func (cc *CurrencyConfig) ToGameConfig(gameName string, gameID uuid.UUID) *CurrencyGameConfig {
//		return &CurrencyGameConfig{
//			IntegratorName: cc.ProviderIntegratorPair.Integrator.Name,
//			ProviderName:   cc.ProviderIntegratorPair.Provider.Name,
//			GameName:       gameName,
//			GameID:         gameID,
//
//			DefaultWager: cc.DefaultWager,
//			WagerLevels:  cc.WagerLevels,
//			Multipliers:  cc.Multipliers,
//		}
//	}
//
//	func (CurrencyConfig) TableName() string {
//		return "currency_configs"
//	}
type Currency struct {
	Title        string  `json:"title"  gorm:"primaryKey"`
	Alias        string  `json:"alias"`
	Type         string  `json:"type"`
	Rate         float64 `json:"additional_alias_rate" gorm:"column:additional_alias_rate"`
	BaseCurrency string  `json:"base_currency"`
}

type CurrencyInfo struct {
	Table      [][]string
	Integrator string
	Provider   string
}

type CurrencyAttributes struct {
	Title      string
	Multiplier int64
	Synonym    string
}

type CurrencyMultiplier struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	OrganizationPairID     uuid.UUID               `json:"organization_pair_id"`
	ProviderIntegratorPair *ProviderIntegratorPair `json:"provider_integrator_pair" gorm:"foreignKey:organization_pair_id"`

	Title      string `json:"title"  gorm:"primaryKey"`
	Multiplier int64  `json:"multiplier"`
	Synonym    string `json:"synonym"`
}

type GroupedCurrencyMultiplier struct {
	ProviderIntegratorPair *ProviderIntegratorPair

	Multipliers map[string]int64
	Synonyms    map[string]string
}

func GroupCurrencyMultiplier(muls []*CurrencyMultiplier) []*GroupedCurrencyMultiplier {
	gr := lo.GroupBy(muls, func(item *CurrencyMultiplier) uuid.UUID {
		return item.OrganizationPairID
	})

	res := []*GroupedCurrencyMultiplier{}

	for _, grouped := range gr {
		resItem := &GroupedCurrencyMultiplier{
			ProviderIntegratorPair: grouped[0].ProviderIntegratorPair,
			Multipliers:            map[string]int64{},
			Synonyms:               map[string]string{},
		}

		for _, item := range grouped {
			resItem.Multipliers[item.Title] = item.Multiplier
			resItem.Synonyms[item.Title] = item.Synonym
		}

		res = append(res, resItem)
	}

	return res
}

func (CurrencyMultiplier) TableName() string {
	return "currency_multipliers"
}

type CurrencyGameConfig struct {
	IntegratorName string    `json:"integrator_name"`
	ProviderName   string    `json:"provider_name"`
	GameName       string    `json:"game_name"`
	GameID         uuid.UUID `json:"game_id"`

	DefaultWager        int64             `json:"default_wager"`
	WagerLevels         []int64           `json:"wager_levels"`
	Multipliers         map[string]int64  `json:"multipliers" gorm:"-"`
	AvailableRTP        []int64           `json:"available_rtp"`
	AvailableVolatility []string          `json:"available_volatility"`
	OnlineVolatility    bool              `json:"online_volatility"`
	GambleDoubleUp      int64             `json:"gamble_double_up"`
	Synonyms            map[string]string `json:"synonyms" gorm:"-"`
}

type CurrencyExchange struct {
	CreatedAt time.Time `json:"created_at"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Rate      float64   `json:"rate"`
}

type CurrencyExchangePagination struct {
	Order    string `json:"order" form:"order"`
	Limit    int    `json:"limit" form:"limit" validate:"required"`
	Page     int    `json:"page" form:"page" validate:"required"`
	Currency string `json:"currency"`
}

func CurrencyFromExchange(out *exchange.CurrencyRates) *CurrencyExchange {
	return &CurrencyExchange{
		CreatedAt: out.CreatedAt.AsTime(),
		From:      out.From,
		To:        out.To,
		Rate:      out.Rate,
	}
}
