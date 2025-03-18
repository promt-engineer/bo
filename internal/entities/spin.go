package entities

import (
	"backoffice/pkg/history"
	"backoffice/utils"
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

const FinanceDivider = 1000
const FloatPrecision = 4

func prettifyPrecision(f float64) float64 {
	return math.Round(f/FinanceDivider*100) / 100
}

type Spin struct {
	CreatedAt time.Time `json:"created_at" csv:"created_at" xlsx:"Created At"`
	UpdatedAt time.Time `json:"updated_at" csv:"updated_at" xlsx:"Updated At"`

	Country   string          `json:"country" gorm:"<-:create"`
	Host      string          `json:"host" csv:"host" xlsx:"Host" gorm:"<-:create"`
	ClientIP  string          `json:"client_ip" csv:"client_ip" gorm:"<-:create"`
	UserAgent string          `json:"user_agent" csv:"user_agent" gorm:"<-:create"`
	Request   json.RawMessage `json:"request" csv:"request" swaggertype:"string" gorm:"<-:create"`

	ID           uuid.UUID `json:"id" gorm:"primaryKey" csv:"id" xlsx:"ID"`
	GameID       uuid.UUID `json:"game_id" gorm:"primaryKey" csv:"game_id" xlsx:"-"`
	Game         string    `json:"game" csv:"game" xlsx:"Game Name"`
	SessionToken uuid.UUID `json:"session_token" csv:"session_token" xlsx:"Session Token"`
	Integrator   string    `json:"integrator" csv:"integrator" xlsx:"Integrator"`
	Operator     string    `json:"operator" csv:"operator" xlsx:"Operator"`

	UserID         string `json:"user_id" csv:"user_id" xlsx:"User ID"`
	ExternalUserID string `json:"external_user_id" csv:"external_user_id" xlsx:"External User ID"`
	Currency       string `json:"currency" csv:"currency" xlsx:"Currency"`

	StartBalance float64 `json:"start_balance" csv:"start_balance" xlsx:"Start Balance"`
	EndBalance   float64 `json:"end_balance" csv:"end_balance" xlsx:"End Balance"`
	Wager        float64 `json:"wager" csv:"wager" xlsx:"Wager"`
	BaseAward    float64 `json:"base_award" csv:"base_award" xlsx:"Base Award"`
	BonusAward   float64 `json:"bonus_award" csv:"bonus_award" xlsx:"Bonus Award"`
	FinalAward   float64 `json:"final_award" csv:"final_award" xlsx:"Final Award"`

	Spin             json.RawMessage `json:"spin" gorm:"serializer:json" csv:"-" swaggertype:"string" xlsx:"-"`
	RestoringIndexes json.RawMessage `json:"restoring_indexes" gorm:"serializer:json" csv:"-" swaggertype:"string" xlsx:"-"`

	IsShown *bool `json:"is_shown" csv:"-" xlsx:"isShown"`
	IsPFR   *bool `json:"is_pfr" csv:"is_pfr" xlsx:"isPFR"`
}

type GameParser struct {
	Game string `json:"game"`
}

func (s *Spin) Prettify() *Spin {
	s.StartBalance = prettifyPrecision(s.StartBalance)
	s.EndBalance = prettifyPrecision(s.EndBalance)
	s.Wager = prettifyPrecision(s.Wager)
	s.BaseAward = prettifyPrecision(s.BaseAward)
	s.BonusAward = prettifyPrecision(s.BonusAward)
	s.FinalAward = prettifyPrecision(s.FinalAward)

	return s
}

func (s *Spin) ToXLSX() utils.XLSXView {
	return s.Prettify()
}

func SpinFromHistory(out *history.SpinOut) *Spin {
	return &Spin{
		CreatedAt:        out.CreatedAt.AsTime(),
		UpdatedAt:        out.UpdatedAt.AsTime(),
		Country:          out.Country,
		Host:             out.Host,
		ClientIP:         out.ClientIp,
		UserAgent:        out.UserAgent,
		Request:          out.Request,
		Operator:         out.Operator,
		ID:               uuid.MustParse(out.Id),
		GameID:           uuid.MustParse(out.GameId),
		Game:             out.Game,
		SessionToken:     uuid.MustParse(out.SessionToken),
		Integrator:       out.Integrator,
		UserID:           out.InternalUserId,
		ExternalUserID:   out.ExternalUserId,
		Currency:         out.Currency,
		StartBalance:     float64(out.StartBalance),
		EndBalance:       float64(out.EndBalance),
		Wager:            float64(out.Wager),
		BaseAward:        float64(out.BaseAward),
		BonusAward:       float64(out.BonusAward),
		FinalAward:       float64(out.FinalAward),
		Spin:             out.Details,
		RestoringIndexes: out.RestoringIndexes,
		IsShown:          out.IsShown,
		IsPFR:            out.IsPfr,
	}
}

type GroupedSpin struct {
	CreatedAt time.Time `json:"created_at" csv:"created_at" xlsx:"Created At" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" csv:"updated_at" xlsx:"Updated At" mapstructure:"updated_at"`

	GameID uuid.UUID `json:"game_id" gorm:"primaryKey" csv:"game_id" xlsx:"-" mapstructure:"game_id"`
	Game   string    `json:"game" csv:"game" xlsx:"Game Name" mapstructure:"game"`

	SessionToken   uuid.UUID `json:"session_token" csv:"session_token" xlsx:"Session Token" mapstructure:"session_token"`
	Integrator     string    `json:"integrator" csv:"integrator" xlsx:"Integrator" mapstructure:"integrator"`
	Operator       string    `json:"operator" csv:"operator" xlsx:"Operator" mapstructure:"operator"`
	Host           *string   `json:"host" csv:"host" xlsx:"Host" mapstructure:"host"`
	UserID         *string   `json:"user_id" csv:"user_id" xlsx:"User ID" mapstructure:"user_id"`
	ExternalUserID *string   `json:"external_user_id" csv:"external_user_id" xlsx:"External User ID" mapstructure:"external_user_id"`
	Currency       string    `json:"currency" csv:"currency" xlsx:"Currency" mapstructure:"currency"`

	Wager      float64 `json:"wager" csv:"wager" xlsx:"Wager" mapstructure:"wager"`
	BaseAward  float64 `json:"base_award" csv:"base_award" xlsx:"Base Award" mapstructure:"base_award"`
	BonusAward float64 `json:"bonus_award" csv:"bonus_award" xlsx:"Bonus Award" mapstructure:"bonus_award"`
	FinalAward float64 `json:"final_award" csv:"final_award" xlsx:"Final Award" mapstructure:"final_award"`

	IsShown *bool `json:"is_shown" csv:"-" xlsx:"isShown" mapstructure:"is_shown"`
	IsPFR   *bool `json:"is_pfr" csv:"is_pfr" xlsx:"isPFR"  mapstructure:"is_pfr"`
}

func (s *GroupedSpin) Prettify() *GroupedSpin {
	s.Wager = prettifyPrecision(s.Wager)
	s.BaseAward = prettifyPrecision(s.BaseAward)
	s.BonusAward = prettifyPrecision(s.BonusAward)
	s.FinalAward = prettifyPrecision(s.FinalAward)

	return s
}

func (*GroupedSpin) TableName() string {
	return "spins"
}
