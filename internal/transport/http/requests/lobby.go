package requests

import "github.com/google/uuid"

type LobbyRequest struct {
	ShortLink    *bool     `json:"short_link"`
	ShowCheats   *bool     `json:"show_cheats"`
	Currency     string    `json:"currency"`
	Game         string    `json:"game"`
	UserID       uuid.UUID `json:"user_id"`
	SessionID    uuid.UUID `json:"session_id"`
	Jurisdiction string    `json:"jurisdiction"`
	UserLocale   string    `json:"user_locale"`
	Integrator   string    `json:"integrator"`
	LobbyURL     string    `json:"lobby_url"`
	RTP          *int64    `json:"rtp"`
	WagerSetID   uuid.UUID `json:"wager_set_id"`
	Volatility   *string   `json:"volatility"`
	LowBalance   *bool     `json:"low_balance"`
}
