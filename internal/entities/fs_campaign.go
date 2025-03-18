package entities

import (
	"github.com/google/uuid"
	"time"
)

type FreeSpinCampaign struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	ID          uuid.UUID
	CreatedBy   uuid.UUID
	Name        string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Type        string

	Token      string
	FSType     string
	Currencies []string
	Status     string
	CoinSize   int
}

func (e FreeSpinCampaign) TableName() string {
	return "fs_campaigns"
}
