package requests

import "github.com/google/uuid"

type UpsertOrganizationRequest struct {
	Name   string `json:"name" validate:"required"`
	Type   string `json:"type" validate:"required"`
	Status int    `json:"status"`
}

type IntegratorGameRequest struct {
	GameID     []uuid.UUID `json:"game_id" validate:"required"`
	WagerSetID uuid.UUID   `json:"wager_set_id"`
}

type IntegratorOneGameRequest struct {
	GameID     uuid.UUID `json:"game_id" validate:"required"`
	WagerSetID uuid.UUID `json:"wager_set_id"`
}

type UpdateIntegratorGameRequest struct {
	IntegratorOneGameRequest
	RTP        *int64  `json:"rtp"`
	Volatility *string `json:"volatility"`
	ShortLink  bool    `json:"short_link"`
}

type IntegratorOperatorPair struct {
	IntegratorID uuid.UUID `json:"integrator_id" form:"integrator_id"`
	OperatorID   uuid.UUID `json:"operator_id" form:"operator_id"`
}

type IntegratorGameWagerSetRequest struct {
	IntegratorOneGameRequest
	Currency string `json:"currency"`
}

type UpdateIntegratorGameWagerSetRequest struct {
	IntegratorOneGameRequest
	Currency    string `json:"currency"`
	NewCurrency string `json:"new_currency"`
}
