package entities

import (
	"backoffice/pkg/history"
	"backoffice/utils"
	"github.com/mitchellh/mapstructure"
)

type FinancialFilters struct {
	Integrator     string `json:"integrator,omitempty" form:"integrator" mapstructure:"integrator,omitempty"`
	Operator       string `json:"operator,omitempty" form:"operator" mapstructure:"operator,omitempty"`
	GameName       string `json:"game,omitempty" form:"game" mapstructure:"game,omitempty"`
	StartingFrom   string `json:"starting_from,omitempty" form:"starting_from" validate:"custom_datetime" mapstructure:"starting_from,omitempty"`
	EndingAt       string `json:"ending_at,omitempty" form:"ending_at" validate:"custom_datetime" mapstructure:"ending_at,omitempty"`
	SessionID      string `json:"session_id" form:"session_id" mapstructure:"session_id,omitempty"`
	RoundID        string `json:"round_id" form:"round_id" mapstructure:"round_id,omitempty"`
	Host           string `json:"host" form:"host" mapstructure:"host,omitempty"`
	UserID         string `json:"user_id" form:"user_id" mapstructure:"user_id,omitempty"`
	ExternalUserID string `json:"external_user_id" form:"external_user_id" mapstructure:"external_user_id,omitempty"`
	IsDemo         *bool  `json:"is_demo,omitempty" form:"is_demo"`
	RTPFrom        int64  `json:"rtp_from,omitempty" form:"rtp_from" mapstructure:"rtp_from,omitempty"`
}

type AggregateFilters struct {
	Currency *string `json:"currency" form:"currency" validate:"required" mapstructure:"-"`

	Integrator   string `json:"integrator,omitempty" form:"integrator" mapstructure:"integrator,omitempty"`
	Operator     string `json:"operator,omitempty" form:"operator" mapstructure:"operator,omitempty"`
	StartingFrom string `json:"starting_from,omitempty" form:"starting_from" validate:"custom_datetime" mapstructure:"starting_from,omitempty"`
	EndingAt     string `json:"ending_at,omitempty" form:"ending_at" validate:"custom_datetime" mapstructure:"ending_at,omitempty"`
	IsPFR        *bool  `json:"is_pfr,omitempty"  form:"is_demo"`
	IsDemo       *bool  `json:"is_demo,omitempty"  form:"is_demo"`
}

type FinancialBase struct {
	Currency *string `json:"currency" form:"currency" validate:"required" mapstructure:"-"`
	FinancialFilters
}

type SpinPagination struct {
	FinancialBase
	Order   string  `json:"order" form:"order"`
	Limit   int     `json:"limit" form:"limit" validate:"required"`
	Page    int     `json:"page" form:"page" validate:"required"`
	GroupBy *string `json:"group_by" form:"group_by"`
}

func (r FinancialFilters) ToSpinFilters() (map[string]interface{}, error) {
	filtersMap := map[string]interface{}{}

	if err := mapstructure.Decode(&r, &filtersMap); err != nil {
		return nil, err
	}

	return filtersMap, nil
}

func (r AggregateFilters) ToSpinFilters() (map[string]interface{}, error) {
	filtersMap := map[string]interface{}{}

	if err := mapstructure.Decode(&r, &filtersMap); err != nil {
		return nil, err
	}

	return filtersMap, nil
}

func (fb *FinancialBase) ToHistoryFilters(gameIDs []string) (*history.FinancialBase, error) {
	start, err := utils.ParseTimestampPB(fb.StartingFrom)
	if err != nil {
		return nil, err
	}

	end, err := utils.ParseTimestampPB(fb.EndingAt)
	if err != nil {
		return nil, err
	}

	hfb := &history.FinancialBase{
		Games: gameIDs,
		Filters: &history.Filters{
			Integrator:     fb.Integrator,
			Operator:       fb.Operator,
			Game:           fb.GameName,
			StartingFrom:   start,
			EndingAt:       end,
			SessionToken:   fb.SessionID,
			RoundId:        fb.RoundID,
			Host:           fb.Host,
			ExternalUserId: fb.ExternalUserID,
			IsDemo:         fb.IsDemo,
			RtpFrom:        fb.RTPFrom,
		},
	}

	if fb.Currency != nil {
		hfb.ConvertCurrency = *fb.Currency
	}

	return hfb, nil
}

func (af *AggregateFilters) ToHistoryFilters(gameIDs []string) (*history.GetAggregatedReportFilters, error) {
	start, err := utils.ParseTimestampPB(af.StartingFrom)
	if err != nil {
		return nil, err
	}
	end, err := utils.ParseTimestampPB(af.EndingAt)
	if err != nil {
		return nil, err
	}

	fil := &history.GetAggregatedReportFilters{
		Games:        gameIDs,
		Integrator:   af.Integrator,
		Operator:     af.Operator,
		StartingFrom: start,
		EndingAt:     end,
		IsPfr:        af.IsPFR,
		IsDemo:       af.IsDemo,
	}

	if af.Currency != nil {
		fil.ConvertCurrency = *af.Currency
	}

	return fil, nil
}
