package entities

import (
	"github.com/google/uuid"
	"time"
)

type AccountRole struct {
	CreatedAt time.Time `json:"created_at"`
	AccountID uuid.UUID `json:"account_id"`
	RoleID    uuid.UUID `json:"role_id"`
}
