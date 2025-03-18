package entities

import (
	"backoffice/pkg/totp"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"time"
)

const (
	AuthProviderEmail = "email"
)

type Account struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	ID uuid.UUID `json:"id"`

	AuthProvider      string `json:"auth_provider"`
	AuthProviderID    string `json:"auth_provider_id"`
	AuthProviderToken string `json:"-"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`

	Organizations []*Organization `json:"organizations" gorm:"many2many:account_organizations;foreignKey:id;joinForeignKey:account_id;joinReferences:organization_id;references:id"`
	Roles         []*Role         `json:"roles" gorm:"many2many:account_roles;foreignKey:id;joinForeignKey:account_id;joinReferences:role_id;references:id"`
	//Permissions []*Permission `json:"permissions" gorm:"many2many:account_permissions;foreignKey:id;joinForeignKey:account_id;joinReferences:permission_id;references:id"`
	// computed
	Permissions []*Permission `json:"permissions" gorm:"-"`

	Status int64 `json:"status"`
	*totp.Params

	ResetPasswordToken     string     `json:"reset_password_token" gorm:"column:reset_password_token"`
	ResetPasswordExpiresAt *time.Time `json:"reset_password_expires_at" gorm:"column:reset_password_expires_at"`

	Operators []*Organization `json:"operators" gorm:"many2many:account_operators;foreignKey:id;joinForeignKey:account_id;joinReferences:operator_id;references:id"`
}

func (a *Account) Compute() {
	for _, role := range a.Roles {
		a.Permissions = append(a.Permissions, role.Permissions...)
	}

	a.Permissions = lo.UniqBy(a.Permissions, func(item *Permission) uuid.UUID {
		return item.ID
	})
}

func (a *Account) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

func (a *Account) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &a)
}

func (a *Account) Authorized(endpoint, method string) bool {
	for _, role := range a.Roles {
		if role.Type == RootRoleTypeName {
			return true
		}

		for _, permission := range role.Permissions {
			if permission.Endpoint == endpoint && permission.IsActionMatched(method) {
				return true
			}
		}
	}

	for _, permission := range a.Permissions {
		if permission.Endpoint == endpoint && permission.IsActionMatched(method) {
			return true
		}
	}

	return false
}

func (a *Account) IsRoot() bool {
	for _, role := range a.Roles {
		if role.Type == RootRoleTypeName {
			return true
		}
	}

	return false
}

func (a *Account) GetDefaultOrganizationID() uuid.UUID {
	if a.Organizations != nil {
		return a.Organizations[0].ID
	}

	return uuid.Nil
}
