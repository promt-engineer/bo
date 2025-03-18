package entities

import "github.com/google/uuid"

type IntegratorOperatorPair struct {
	ID           uuid.UUID `json:"id"`
	IntegratorID uuid.UUID `json:"integrator_id"`
	OperatorID   uuid.UUID `json:"operator_id"`

	Integrator *Organization `json:"integrator,omitempty" gorm:"foreignKey:IntegratorID"`
	Operator   *Organization `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

func (IntegratorOperatorPair) TableName() string {
	return "operator_integrators"
}
