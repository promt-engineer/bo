package entities

import (
	"github.com/google/uuid"
	"time"
)

type AccountPermission struct {
	CreatedAt    time.Time `json:"created_at"`
	AccountID    uuid.UUID `json:"account_id"`
	PermissionID uuid.UUID `json:"permission_id"`
}
