package entities

import (
	"backoffice/pkg/history"

	"github.com/google/uuid"
)

type SimpleAggregatedReport struct {
	Game   string    `json:"game" csv:"game" xlsx:"Game" `
	GameID uuid.UUID `json:"game_id" csv:"game_id" xlsx:"GameID"`

	Currency string `json:"currency" csv:"currency" xlsx:"Currency"`

	UserCount  int `json:"user_count" csv:"user_count" xlsx:"UserCount"`
	RoundCount int `json:"round_count" csv:"round_count" xlsx:"RoundCount"`

	Wager float64 `json:"wager" csv:"wager" xlsx:"Wager"`
	Award float64 `json:"award" csv:"award" xlsx:"Award"`
}

type AggregatedReportByGame struct {
	Game   string    `json:"game" csv:"game" xlsx:"Game"`
	GameID uuid.UUID `json:"game_id" csv:"game_id" xlsx:"GameID"`

	AggregatedReport
}

func AggregatedReportByGameFromHistory(out *history.GetAggregatedReportByGameItem) *AggregatedReportByGame {
	return &AggregatedReportByGame{
		Game:   out.Game,
		GameID: uuid.MustParse(out.GameId),

		AggregatedReport: AggregatedReportFromHistory(out),
	}
}

type AggregatedReportByCountry struct {
	Country string `json:"country" csv:"country" xlsx:"Country"`

	AggregatedReport
}

func AggregatedReportByCountryFromHistory(out *history.GetAggregatedReportByCountryItem) *AggregatedReportByCountry {
	return &AggregatedReportByCountry{
		Country: out.Country,

		AggregatedReport: AggregatedReportFromHistory(out),
	}
}

type AggregatedReport struct {
	Currency string `json:"currency" csv:"currency" xlsx:"Currency"`

	UserCount     int `json:"user_count" csv:"user_count" xlsx:"UserCount"`
	RoundCount    int `json:"round_count" csv:"round_count" xlsx:"RoundCount"`
	PFRRoundCount int `json:"pfr_round_count" csv:"pfr_round_count" xlsx:"PFRRoundCount"`

	Wager float64 `json:"wager" csv:"wager" xlsx:"Wager"`
	Award float64 `json:"award" csv:"award" xlsx:"Award"`

	PFRWager float64 `json:"pfr_wager" csv:"pfr_wager" xlsx:"PFRWager"`
	PFRAward float64 `json:"pfr_award" csv:"pfr_award" xlsx:"PFRAward"`

	// computed
	RTP float64 `json:"rtp" csv:"rtp" xlsx:"RTP"`

	Revenue    float64 `json:"revenue" csv:"revenue" xlsx:"Revenue"`
	PFRRevenue float64 `json:"pfr_revenue" csv:"pfr_revenue" xlsx:"PFRRevenue"`

	RoundPerUser   float64 `json:"round_per_user" csv:"round_per_user" xlsx:"RoundPerUser"`
	WagerPerUser   float64 `json:"wager_per_user" csv:"wager_per_user" xlsx:"WagerPerUser"`
	AwardPerUser   float64 `json:"award_per_user" csv:"award_per_user" xlsx:"AwardPerUser"`
	RevenuePerUser float64 `json:"revenue_per_user" csv:"revenue_per_user" xlsx:"RevenuePerUser"`
}

func NewAggregatedReport(currency string, userC, roundC int, wager, award float64) AggregatedReport {
	return AggregatedReport{
		Currency: currency,

		UserCount:  userC,
		RoundCount: roundC,

		Wager: wager,
		Award: award,
	}
}

type HistoryAggregatedReport interface {
	GetCurrency() string

	GetUserCount() int64
	GetRoundCount() int64

	GetWager() float64
	GetAward() float64
}

func AggregatedReportFromHistory(out HistoryAggregatedReport) AggregatedReport {
	return NewAggregatedReport(
		out.GetCurrency(), int(out.GetUserCount()), int(out.GetRoundCount()), out.GetWager(), out.GetAward(),
	)
}

func (ar *AggregatedReport) ExchangeCurrencyMut(toCurrency string, multiplier float64) {
	ar.Currency = toCurrency

	ar.Wager = ar.Wager * multiplier
	ar.Award = ar.Award * multiplier
	ar.PFRWager = ar.PFRWager * multiplier
	ar.PFRAward = ar.PFRAward * multiplier
}

func (ar *AggregatedReport) GetCurrency() string {
	return ar.Currency
}

func (ar *AggregatedReport) Compute() {
	ar.Revenue = ar.Wager - ar.Award
	ar.PFRRevenue = ar.PFRWager - ar.PFRAward

	if ar.UserCount > 0 {
		ar.RoundPerUser = float64(ar.RoundCount) / float64(ar.UserCount)
	}
	if ar.UserCount > 0 {
		ar.WagerPerUser = float64(ar.Wager) / float64(ar.UserCount)
	}

	if ar.UserCount > 0 {
		ar.AwardPerUser = float64(ar.Award) / float64(ar.UserCount)
	}

	if ar.UserCount > 0 {
		ar.RevenuePerUser = float64(ar.Revenue) / float64(ar.UserCount)
	}

	if ar.Wager > 0 {
		ar.RTP = float64(ar.Award) / float64(ar.Wager)
	}
}

func (r *AggregatedReport) Prettify() {
	r.Compute()

	r.Wager = prettifyPrecision(r.Wager)
	r.Award = prettifyPrecision(r.Award)

	r.PFRWager = prettifyPrecision(r.PFRWager)
	r.PFRAward = prettifyPrecision(r.PFRAward)

	r.Revenue = prettifyPrecision(r.Revenue)
	r.PFRRevenue = prettifyPrecision(r.PFRRevenue)

	r.WagerPerUser = prettifyPrecision(r.WagerPerUser)
	r.AwardPerUser = prettifyPrecision(r.AwardPerUser)
	r.RevenuePerUser = prettifyPrecision(r.RevenuePerUser)
}
