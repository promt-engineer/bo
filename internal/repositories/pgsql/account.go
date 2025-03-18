package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accountRepository struct {
	conn *gorm.DB
}

func NewAccountRepository(conn *gorm.DB) *accountRepository {
	return &accountRepository{
		conn: conn,
	}
}

func (r *accountRepository) FindBy(ctx context.Context, params map[string]interface{}) (account *entities.Account, err error) {
	if err = r.conn.WithContext(ctx).Where(params).
		Preload("Roles.Permissions").
		//Preload("Permissions").
		Preload("Organizations").
		Preload("Operators").
		First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	account.Compute()

	return
}

func (r *accountRepository) FindAllByRoleID(ctx context.Context, roleID string) (accounts []*entities.Account, err error) {
	if err = r.conn.WithContext(ctx).Preload("Roles").Preload("Organizations").
		Joins("join account_roles on account_roles.account_id = accounts.id").
		Where("account_roles.role_id = ?", roleID).Find(&accounts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	return
}

func (r *accountRepository) Paginate(ctx context.Context, organizationID uuid.UUID, filters repositories.Filters, order string, limit int, offset int) (accounts []*entities.Account, total int64, err error) {
	query := r.conn.WithContext(ctx).Model(&entities.Account{}).
		Where(filters.Where).
		Not(filters.Not).
		Preload("Organizations").
		Preload("Roles.Permissions").
		Preload("Operators")

	query = query.Joins("join account_organizations on accounts.id = account_organizations.account_id and organization_id = ?", organizationID)

	query = query.Joins("left join account_operators on accounts.id = account_operators.account_id")

	if err = query.Count(&total).Error; err != nil {
		return
	}

	if err = query.Order(order).Limit(limit).Offset(offset).Find(&accounts).Error; err != nil {
		return
	}
	return
}

func (r *accountRepository) Create(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	if err := r.conn.WithContext(ctx).Create(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (r *accountRepository) Delete(ctx context.Context, account *entities.Account) error {
	return r.conn.WithContext(ctx).Where("id = ?", account.ID).Delete(&account).Error
}

func (r *accountRepository) Save(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	if err := r.conn.Save(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (r *accountRepository) Update(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	if err := r.conn.WithContext(ctx).Model(&account).Updates(account).Error; err != nil {
		return nil, err
	}

	return r.FindBy(ctx, map[string]interface{}{"id": account.ID})
}

func (r *accountRepository) DisableTOTP(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	if err := r.conn.Omit(clause.Associations).WithContext(ctx).Model(&account).Updates(map[string]interface{}{
		"totp_enabled": false,
		"totp_secret":  "",
		"totp_url":     "",
	}).Error; err != nil {
		return nil, err
	}

	return r.FindBy(ctx, map[string]interface{}{"id": account.ID})
}
