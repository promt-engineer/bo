package entities

import (
	"backoffice/utils"
	"sort"
	"time"
)

type UserReport struct {
	CreatedAt time.Time `json:"created_at" csv:"created_at" xlsx:"Created At"`
	UpdatedAt time.Time `json:"updated_at" csv:"updated_at" xlsx:"Updated At"`

	Wager           float64 `json:"wager" csv:"wager" xlsx:"Wager"`
	WagerWithoutPFR float64 `json:"wager_without_pfr" csv:"wager_without_pfr" xlsx:"Wager without PFR"`
	Award           float64 `json:"award" csv:"award" xlsx:"Award"`
	AwardWithoutPFR float64 `json:"award_without_pfr" csv:"award_without_pfr" xlsx:"Award without PFR"`

	Currency   string `json:"currency" xlsx:"Currency"`
	Integrator string `json:"integrator" csv:"integrator" xlsx:"Integrator"`
	Operator   string `json:"operator" csv:"operator" xlsx:"Operator"`

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
}

func (gs *UserReport) TableName() string {
	return "spins"
}

func (gs *UserReport) Compute() *UserReport {
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

	return gs
}

func (gs *UserReport) Prettify() *UserReport {
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

	return gs
}

func (s *UserReport) GetIntegrator() string {
	return s.Integrator
}

func (s *UserReport) GetCurrency() string {
	return s.Currency
}

func (gs *UserReport) ToXLSX() utils.XLSXView {
	return gs.Prettify()
}

func UserReportFromSpins(out []*Spin) *UserReport {
	if len(out) == 0 {
		return &UserReport{}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})

	sb := out[0].StartBalance
	eb := out[len(out)-1].EndBalance

	var wager float64
	var wagerWithoutPfr float64
	var award float64
	var awardWithoutPfr float64

	for _, out := range out {
		totalAward := out.FinalAward

		if out.IsPFR != nil && !*out.IsPFR {
			awardWithoutPfr += totalAward
			wagerWithoutPfr += out.Wager
		}
		wager += out.Wager
		award += totalAward

	}

	return &UserReport{
		CreatedAt:       out[0].CreatedAt,
		UpdatedAt:       out[len(out)-1].CreatedAt,
		Wager:           wager,
		WagerWithoutPFR: wagerWithoutPfr,
		Award:           award,
		AwardWithoutPFR: awardWithoutPfr,
		Currency:        out[0].Currency,
		Integrator:      out[0].Integrator,
		Operator:        out[0].Operator,
		StartBalance:    &sb,
		EndBalance:      &eb,
	}
}
