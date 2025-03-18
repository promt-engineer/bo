package entities

import "github.com/google/uuid"

type ProviderIntegratorPair struct {
	ID           uuid.UUID `json:"id"`
	ProviderID   uuid.UUID `json:"provider_id"`
	IntegratorID uuid.UUID `json:"integrator_id"`

	Provider   *Organization `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	Integrator *Organization `json:"integrator,omitempty" gorm:"foreignKey:IntegratorID"`
}

func (ProviderIntegratorPair) TableName() string {
	return "integrator_providers"
}
