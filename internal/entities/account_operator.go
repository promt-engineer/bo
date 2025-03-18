package entities

import (
	"github.com/google/uuid"
	"time"
)

type AccountOperator struct {
	CreatedAt  time.Time `json:"created_at"`
	AccountID  uuid.UUID `json:"account_id"`
	OperatorID uuid.UUID `json:"operator_id"`
}
