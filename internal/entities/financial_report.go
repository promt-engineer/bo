package entities

import (
	"backoffice/pkg/history"
	"backoffice/utils"
	"fmt"
)

type FinancialReport struct {
	Award           float64 `json:"award" csv:"award" xlsx:"Award"`
	AwardWithoutPFR float64 `json:"award_without_pfr" csv:"award_without_pfr" xlsx:"Award without PFR"`
	Wager           float64 `json:"wager" csv:"wager" xlsx:"Wager"`
	WagerWithoutPFR float64 `json:"wager_without_pfr" csv:"wager_without_pfr" xlsx:"Wager without PFR"`

	SpinQuantity int `json:"spin_quantity" csv:"spin_quantity" xlsx:"Spin Quantity"`
	UserQuantity int `json:"user_quantity" csv:"user_quantity" xlsx:"User Quantity"`

	// computed
	Revenue         float64 `json:"revenue" csv:"revenue" xlsx:"Revenue"`
	RTP             float64 `json:"rtp" csv:"rtp" xlsx:"RTP"`
	RTPWithTurnover float64 `json:"rtp_with_turnover" csv:"rtp_with_turnover" xlsx:"RTP With Turnover"`

	Margin       float64 `json:"margin"  csv:"margin" xlsx:"Margin"`
	AwardWithPFR float64 `json:"award_with_pfr" csv:"award_with_pfr" xlsx:"Award with PFR"`
	WagerWithPFR float64 `json:"wager_with_pfr" csv:"wager_with_pfr" xlsx:"Wager with PFR"`
}

func (r *FinancialReport) Prettify() *FinancialReport {
	if r.Wager != 0 {
		r.RTP = r.Award / r.Wager
	}

	if r.WagerWithoutPFR != 0 {
		r.RTPWithTurnover = r.Award / r.WagerWithoutPFR
		r.Margin = 1 - r.RTPWithTurnover
	}

	r.Revenue = r.WagerWithoutPFR - r.Award

	r.Revenue = prettifyPrecision(r.Revenue)
	r.AwardWithPFR = prettifyPrecision(r.Award - r.AwardWithoutPFR)
	r.WagerWithPFR = prettifyPrecision(r.Wager - r.WagerWithoutPFR)

	r.Award = prettifyPrecision(r.Award)
	r.Wager = prettifyPrecision(r.Wager)
	r.AwardWithoutPFR = prettifyPrecision(r.AwardWithoutPFR)
	r.WagerWithoutPFR = prettifyPrecision(r.WagerWithoutPFR)
	r.AwardWithPFR = prettifyPrecision(r.AwardWithPFR)
	r.WagerWithPFR = prettifyPrecision(r.WagerWithPFR)

	return r
}

func FinancialReportFromHistory(out *history.FinancialReport) *FinancialReport {
	return &FinancialReport{
		Award:           float64(out.Award),
		AwardWithoutPFR: float64(out.AwardWithoutPfr),
		Wager:           float64(out.Wager),
		WagerWithoutPFR: float64(out.WagerWithoutPfr),

		SpinQuantity: int(out.SpinQuantity),
		UserQuantity: int(out.UserQuantity),
	}
}

type RawFinancialReport struct {
	Currency        string
	Integrator      string
	Award           float64
	Wager           float64
	AwardWithoutPFR float64
	WagerWithoutPFR float64
}

type RawUserStatReport struct {
	Currency   string
	Integrator string

	SpinQuantity int `json:"spin_quantity"`
	UserQuantity int `json:"user_quantity"`
}

// ToFrontendView immutable func
func (r *FinancialReport) ToXLSX() utils.XLSXView {
	return r.Prettify()
}

type FinancialReportKey struct {
	Currency   string
	Integrator string
}

func (frk *FinancialReportKey) String() string {
	return fmt.Sprintf("cur - %s int - %s", frk.Currency, frk.Integrator)
}
