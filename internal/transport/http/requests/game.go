package requests

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type GameRequest struct {
	Name                 string         `json:"name" validate:"required"`
	Jurisdictions        pq.StringArray `json:"jurisdictions" gorm:"type:varchar[]" swaggertype:"array,string" validate:"required"`
	Currencies           pq.StringArray `json:"currencies" gorm:"type:varchar[]" swaggertype:"array,string" validate:"required"`
	Languages            pq.StringArray `json:"languages" gorm:"type:varchar[]" swaggertype:"array,string" validate:"required"`
	UserLocales          pq.StringArray `json:"user_locales" gorm:"type:varchar[]" swaggertype:"array,string" validate:"required"`
	ApiURL               string         `json:"api_url"`
	ClientURL            string         `json:"client_url"`
	OrganizationID       uuid.UUID      `json:"organization_id" validate:"required"`
	WagerSetID           uuid.UUID      `json:"wager_set_id" validate:"required"`
	IsPublic             *bool          `json:"is_public" validate:"required"`
	IsStatisticShown     *bool          `json:"is_statistic_shown" validate:"required"`
	IsDemo               *bool          `json:"is_demo" validate:"required"`
	IsFreeSpins          *bool          `json:"is_freespins" validate:"required"`
	RTP                  *int64         `json:"rtp"`
	Volatility           *string        `json:"volatility"`
	AvailableRTP         pq.Int64Array  `json:"available_rtp" gorm:"type:integer[]" swaggertype:"array,integer"`
	AvailableVolatility  pq.StringArray `json:"available_volatility" gorm:"type:varchar[]" swaggertype:"array,string"`
	OnlineVolatility     *bool          `json:"online_volatility" validate:"required"`
	AvailableWagerSetsID pq.StringArray `json:"available_wager_sets_id" gorm:"type:uuid[]" swaggertype:"array,string" validate:"required"`
	GambleDoubleUp       int64          `json:"gamble_double_up"`
}

type GameListRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
}
