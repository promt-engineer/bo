package entities

import "github.com/google/uuid"

const (
	RootRoleTypeName  = "root"
	AdminRoleTypeName = "admin"
)

type Role struct {
	ID             uuid.UUID     `json:"id"`
	OrganizationID uuid.UUID     `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Type           string        `json:"type"`
	Permissions    []*Permission `json:"permissions" gorm:"many2many:role_permissions;foreignKey:id;joinForeignKey:role_id;joinReferences:permission_id;references:id"`
}
