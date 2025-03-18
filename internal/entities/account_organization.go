package entities

import (
	"github.com/google/uuid"
	"time"
)

type AccountOrganization struct {
	CreatedAt      time.Time `json:"created_at"`
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
}
