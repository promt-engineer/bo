package repositories

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
)

type RoleRepository interface {
	GetAccountRoles(ctx context.Context, account *entities.Account) ([]*entities.Role, error)
	Paginate(ctx context.Context, organizationID uuid.UUID, filters Filters, order string, limit int, offset int) (roles []*entities.Role, total int64, err error)
	Create(ctx context.Context, role *entities.Role) (*entities.Role, error)
	Update(ctx context.Context, role *entities.Role) (*entities.Role, error)
	Assign(ctx context.Context, account *entities.Account, role *entities.Role) error
	Revoke(ctx context.Context, account *entities.Account, role *entities.Role) error
	Delete(ctx context.Context, role *entities.Role) error
	FindBy(ctx context.Context, params map[string]interface{}) (role *entities.Role, err error)
}
