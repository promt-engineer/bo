package entities

import (
	"github.com/google/uuid"
	"time"
)

const (
	OrganizationTypeIntegrator = "integrator"
	OrganizationTypeProvider   = "provider"
	OrganizationTypeOperator   = "operator"
)

type Organization struct {
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Status *uint8    `json:"status"`
	ApiKey string    `json:"api_key"`
}

func (o *Organization) IsIntegrator() bool {
	return o.Type == OrganizationTypeIntegrator
}

func (o *Organization) IsOperator() bool {
	return o.Type == OrganizationTypeOperator
}
