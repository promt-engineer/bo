package entities

import (
	"backoffice/internal/errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"sort"
	"time"
)

type WagerSet struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ID             uuid.UUID     `json:"id"`
	OrganizationID uuid.UUID     `json:"organization_id"`
	Name           string        `json:"name"`
	WagerLevels    pq.Int64Array `json:"wager_levels" gorm:"type:integer[]" swaggertype:"array,integer"`
	DefaultWager   int64         `json:"default_wager"`

	IsActive bool `json:"is_active"`
}

func (*WagerSet) TableName() string {
	return "wager_sets"
}

func (ws *WagerSet) SetNewWagerParams(wl []int64, dw int64) error {
	if !lo.Contains(wl, dw) {
		return errors.ErrDefaultWagerOutOfList
	}

	if dw <= 0 || lo.ContainsBy(wl, func(item int64) bool {
		return item <= 0
	}) {
		return errors.ErrNegativeWager
	}

	wl = lo.Uniq(wl)

	sort.Slice(wl, func(i, j int) bool {
		return wl[i] < wl[j]
	})

	ws.WagerLevels = wl
	ws.DefaultWager = dw

	return nil
}
