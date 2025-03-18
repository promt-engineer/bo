package repositories

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
)

type AccountRepository interface {
	FindBy(ctx context.Context, params map[string]interface{}) (account *entities.Account, err error)
	FindAllByRoleID(ctx context.Context, roleID string) (accounts []*entities.Account, err error)
	Paginate(ctx context.Context, organizationID uuid.UUID, filters Filters, order string, limit int, offset int) (accounts []*entities.Account, total int64, err error)
	Create(ctx context.Context, account *entities.Account) (*entities.Account, error)
	Delete(ctx context.Context, account *entities.Account) error
	Save(ctx context.Context, account *entities.Account) (*entities.Account, error)
	Update(ctx context.Context, account *entities.Account) (*entities.Account, error)
	DisableTOTP(ctx context.Context, account *entities.Account) (*entities.Account, error)
}
