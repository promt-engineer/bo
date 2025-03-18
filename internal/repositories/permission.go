package repositories

import (
	"backoffice/internal/entities"
	"context"
)

type PermissionRepository interface {
	FindBy(ctx context.Context, params map[string]interface{}) (permission *entities.Permission, err error)
	Create(ctx context.Context, permission *entities.Permission) (*entities.Permission, error)
	All(ctx context.Context) (permissions []*entities.Permission, err error)
	Update(ctx context.Context, permission *entities.Permission) (*entities.Permission, error)
	Delete(ctx context.Context, permission *entities.Permission) error
	GetAccountPermissions(ctx context.Context, account *entities.Account) ([]*entities.Permission, error)
	Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) (permissions []*entities.Permission, total int64, err error)
	RevokeRolePermissions(ctx context.Context, role *entities.Role, permissions ...*entities.Permission) error
	AssignRolePermissions(ctx context.Context, role *entities.Role, permissions []*entities.Permission) error
	//AssignAccountPermissions(ctx context.Context, account *entities.Account, permissions []*entities.Permission) error
	//RevokeAccountPermissions(ctx context.Context, account *entities.Account, permissions ...*entities.Permission) error
}
