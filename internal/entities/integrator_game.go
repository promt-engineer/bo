package entities

import "github.com/google/uuid"

type IntegratorGame struct {
	OrganizationID uuid.UUID `json:"organization_id" xlsx:"-"`
	Organization   *Organization
	GameID         uuid.UUID `json:"game_id" xlsx:"-"`
	Game           *Game
	WagerSetID     uuid.UUID `json:"wager_set_id" xlsx:"-"`
	WagerSet       *WagerSet `json:"wager_set" xlsx:"-" gorm:"foreignKey:WagerSetID;references:ID"`
	RTP            *int64    `json:"rtp" xlsx:"-"`
	Volatility     *string   `json:"volatility" xlsx:"-"`
	ShortLink      bool      `json:"short_link" xlsx:"-"`
}

type IntegratorGameWagerSet struct {
	OrganizationID uuid.UUID `json:"organization_id" xlsx:"-"`
	Organization   *Organization
	GameID         uuid.UUID `json:"game_id" xlsx:"-"`
	Game           *Game
	Currency       string    `json:"currency" xlsx:"-"`
	WagerSetID     uuid.UUID `json:"wager_set_id" xlsx:"-"`
	WagerSet       *WagerSet `json:"wager_set" xlsx:"-" gorm:"foreignKey:WagerSetID;references:ID"`
}
