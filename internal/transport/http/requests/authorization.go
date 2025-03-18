package requests

import "github.com/google/uuid"

type UpsertRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type        string `json:"type" validate:"omitempty"`
}

type AssignAccountRoleRequest struct {
	AccountID uuid.UUID `json:"account_id"`
	RoleID    uuid.UUID `json:"role_id"`
}

type AccountRoleRequest struct {
	RoleID string `json:"role_id" validate:"required"`
}

type AccountOrganizationRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
}

type AccountPermissionsRequest struct {
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

type RolePermissionsRequest struct {
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

type AccountOperatorRequest struct {
	IntegratorID uuid.UUID `json:"integrator_id" validate:"required"`
	OperatorID   uuid.UUID `json:"operator_id" validate:"required"`
}
