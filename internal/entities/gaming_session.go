package entities

import (
	"backoffice/pkg/history"
	"backoffice/utils"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type GamingSession struct {
	CreatedAt time.Time `json:"created_at" csv:"created_at" xlsx:"Created At"`

	SessionToken   uuid.UUID `json:"session_token" csv:"session_token" xlsx:"Session Token"`
	UserID         uuid.UUID `json:"user_id" csv:"user_id" xlsx:"User ID"`
	ExternalUserID string    `json:"external_user_id" csv:"external_user_id" xlsx:"External User ID"`

	Wager           float64 `json:"wager" csv:"wager" xlsx:"Wager"`
	WagerWithoutPFR float64 `json:"wager_without_pfr" csv:"wager_without_pfr" xlsx:"Wager without PFR"`
	Award           float64 `json:"award" csv:"award" xlsx:"Award"`
	AwardWithoutPFR float64 `json:"award_without_pfr" csv:"award_without_pfr" xlsx:"Award without PFR"`

	Currency   string    `json:"currency" xlsx:"Currency"`
	Integrator string    `json:"integrator" csv:"integrator" xlsx:"Integrator"`
	GameID     uuid.UUID `json:"game_id" gorm:"primaryKey" csv:"game_id" xlsx:"-"`
	Game       string    `json:"game" xlsx:"Game Name"`
	Operator   string    `json:"operator" csv:"operator" xlsx:"Operator"`
	Spins      []*Spin   `json:"spins" gorm:"many2many:spin_spins" xlsx:"-"`

	// computed
	WagerWithPFR float64 `json:"wager_with_pfr" csv:"wager_with_pfr" xlsx:"Wager with PFR" gorm:"-"`
	AwardWithPFR float64 `json:"award_with_pfr" csv:"award_with_pfr" xlsx:"Award with PFR" gorm:"-"`

	RTP             float64 `json:"rtp" csv:"rtp" xlsx:"RTP" gorm:"-"`
	RTPWithTurnover float64 `json:"rtp_with_turnover" csv:"rtp_with_turnover" xlsx:"RTP With Turnover" gorm:"-"`
	Margin          float64 `json:"margin" gorm:"-"`
	Revenue         float64 `json:"revenue" csv:"revenue" xlsx:"Revenue" gorm:"-"`

	// spins aggregations
	StartBalance *float64 `json:"start_balance,omitempty" csv:"-" xlsx:"-" gorm:"-"`
	EndBalance   *float64 `json:"end_balance,omitempty" csv:"-" xlsx:"-" gorm:"-"`
	BaseAward    *float64 `json:"base_award,omitempty" csv:"-" xlsx:"-" gorm:"-"`
	BonusAward   *float64 `json:"bonus_award,omitempty" csv:"-" xlsx:"-" gorm:"-"`
	FinalAward   *float64 `json:"final_award,omitempty" csv:"-" xlsx:"-" gorm:"-"`
}

func (gs *GamingSession) TableName() string {
	return "spins"
}

func (gs *GamingSession) Compute() {
	if gs.Wager != 0 {
		gs.RTP = gs.Award / gs.Wager
	}

	if gs.WagerWithoutPFR != 0 {
		gs.RTPWithTurnover = gs.Award / gs.WagerWithoutPFR
	}

	gs.AwardWithPFR = gs.Award - gs.AwardWithoutPFR
	gs.WagerWithPFR = gs.Wager - gs.WagerWithoutPFR

	gs.Margin = 1 - gs.RTPWithTurnover
	gs.Revenue = gs.Wager - gs.Award

	if len(gs.Spins) != 0 {
		sb := gs.Spins[0].StartBalance
		eb := gs.Spins[len(gs.Spins)-1].EndBalance
		baseAward := lo.Reduce(gs.Spins, func(agg float64, item *Spin, index int) float64 {
			return agg + item.BaseAward
		}, 0)

		bonusAward := lo.Reduce(gs.Spins, func(agg float64, item *Spin, index int) float64 {
			return agg + item.BonusAward
		}, 0)

		finalAward := lo.Reduce(gs.Spins, func(agg float64, item *Spin, index int) float64 {
			return agg + item.FinalAward
		}, 0)

		gs.StartBalance = &sb
		gs.EndBalance = &eb
		gs.BaseAward = &baseAward
		gs.BonusAward = &bonusAward
		gs.FinalAward = &finalAward
	}
}

func (gs *GamingSession) Prettify() *GamingSession {
	gs.Compute()
	gs.Wager = prettifyPrecision(gs.Wager)
	gs.WagerWithoutPFR = prettifyPrecision(gs.WagerWithoutPFR)
	gs.WagerWithPFR = prettifyPrecision(gs.WagerWithPFR)
	gs.Award = prettifyPrecision(gs.Award)
	gs.AwardWithoutPFR = prettifyPrecision(gs.AwardWithoutPFR)
	gs.AwardWithPFR = prettifyPrecision(gs.AwardWithPFR)

	gs.Revenue = prettifyPrecision(gs.Revenue)

	if gs.StartBalance != nil {
		*gs.StartBalance = prettifyPrecision(*gs.StartBalance)
	}
	if gs.EndBalance != nil {
		*gs.EndBalance = prettifyPrecision(*gs.EndBalance)
	}
	if gs.BaseAward != nil {
		*gs.BaseAward = prettifyPrecision(*gs.BaseAward)
	}
	if gs.BonusAward != nil {
		*gs.BonusAward = prettifyPrecision(*gs.BonusAward)
	}

	if gs.FinalAward != nil {
		*gs.FinalAward = prettifyPrecision(*gs.FinalAward)
	}

	if len(gs.Spins) != 0 {
		gs.Spins = lo.Map(gs.Spins, func(item *Spin, index int) *Spin {
			return item.Prettify()
		})
	}

	return gs
}

func (gs *GamingSession) ToXLSX() utils.XLSXView {
	return gs.Prettify()
}

func GamingSessionFromHistory(out *history.GameSessionOut) *GamingSession {
	sb := float64(out.StartBalance)
	eb := float64(out.EndBalance)
	baseAward := float64(out.BaseAward)
	bonusAward := float64(out.BonusAward)
	finalAward := float64(out.FinalAward)

	var spins []*Spin

	lo.ForEach(out.Spins, func(item *history.SpinOut, index int) {
		spins = append(spins, SpinFromHistory(item))
	})

	return &GamingSession{
		CreatedAt:       out.CreatedAt.AsTime(),
		SessionToken:    uuid.MustParse(out.SessionToken),
		UserID:          uuid.MustParse(out.UserId),
		ExternalUserID:  out.ExternalUserId,
		Operator:        out.Operator,
		Wager:           float64(out.Wager),
		WagerWithoutPFR: float64(out.WagerWithoutPfr),
		Award:           float64(out.Award),
		AwardWithoutPFR: float64(out.AwardWithoutPfr),
		Currency:        out.Currency,
		Integrator:      out.Integrator,
		GameID:          uuid.MustParse(out.GameId),
		Game:            out.Game,
		Spins:           spins,
		WagerWithPFR:    float64(out.WagerWithPfr),
		AwardWithPFR:    float64(out.AwardWithPfr),
		RTP:             float64(out.Rtp),
		RTPWithTurnover: float64(out.RtpWithTurnover),
		Margin:          float64(out.Margin),
		Revenue:         float64(out.Revenue),
		StartBalance:    &sb,
		EndBalance:      &eb,
		BaseAward:       &baseAward,
		BonusAward:      &bonusAward,
		FinalAward:      &finalAward,
	}
}
