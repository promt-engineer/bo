package entities

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time
	AccountID    uuid.UUID
	AccessToken  string
	RefreshToken string
	ExpiredAt    time.Time
}
