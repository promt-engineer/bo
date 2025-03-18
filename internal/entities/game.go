package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Game struct {
	ID             uuid.UUID      `json:"id" gorm:"primaryKey;column:id" xlsx:"-"`
	OrganizationID uuid.UUID      `json:"organization_id,omitempty" xlsx:"-"`
	Organization   *Organization  `json:"organization,omitempty" xlsx:"-"`
	Name           string         `json:"name" xlsx:"Name"`
	Jurisdictions  pq.StringArray `json:"jurisdictions" gorm:"type:varchar[]" swaggertype:"array,string" xlsx:"-"`
	Currencies     pq.StringArray `json:"currencies" gorm:"type:varchar[]" swaggertype:"array,string" xlsx:"-"`
	Languages      pq.StringArray `json:"languages" gorm:"type:varchar[]" swaggertype:"array,string" xlsx:"-"`
	UserLocales    pq.StringArray `json:"user_locales" gorm:"type:varchar[]" swaggertype:"array,string" xlsx:"-"`
	ApiUrl         string         `json:"api_url" xlsx:"-"`
	ClientUrl      string         `json:"client_url" xlsx:"-"`

	WagerSetID uuid.UUID `json:"wager_set_id" xlsx:"-"`
	WagerSet   *WagerSet `json:"wager_set" xlsx:"-"`

	IsPublic         bool `json:"is_public" xlsx:"-" gorm:"default:false"`
	IsStatisticShown bool `json:"is_statistic_shown" xlsx:"-" gorm:"default:false"`

	IsDemo      bool `json:"is_demo" xlsx:"-" gorm:"default:false"`
	IsFreespins bool `json:"is_freespins" xlsx:"-" gorm:"default:false"`

	RTP        *int64  `json:"rtp" xlsx:"-"`
	Volatility *string `json:"volatility" xlsx:"-"`

	AvailableRTP         pq.Int64Array  `json:"available_rtp" gorm:"type:integer[]" swaggertype:"array,integer" xlsx:"-"`
	AvailableVolatility  pq.StringArray `json:"available_volatility" gorm:"type:varchar[]" swaggertype:"array,string" xlsx:"-"`
	OnlineVolatility     bool           `json:"online_volatility" xlsx:"-" gorm:"default:false"`
	AvailableWagerSetsID pq.StringArray `json:"available_wager_sets_id" gorm:"type:uuid[]" swaggertype:"array,string" xlsx:"-"`
	AvailableWagerSets   []WagerSet     `json:"available_wager_sets" gorm:"-" xlsx:"-"`
	GambleDoubleUp       int64          `json:"gamble_double_up" xlsx:"-"`
}

func (Game) TableName() string {
	return "games"
}

func (Game) JoinTableName(table string) string {
	return "games"
}
