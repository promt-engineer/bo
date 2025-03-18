package entities

import (
	"github.com/google/uuid"
	"time"
)

type RolePermission struct {
	CreatedAt    time.Time `json:"created_at"`
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}
