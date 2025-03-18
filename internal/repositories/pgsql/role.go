package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	conn *gorm.DB
}

func NewRoleRepository(conn *gorm.DB) *roleRepository {
	return &roleRepository{
		conn: conn,
	}
}

func (r *roleRepository) Assign(ctx context.Context, account *entities.Account, role *entities.Role) error {
	var ar *entities.AccountRole

	err := r.conn.WithContext(ctx).Where("account_id = ? and role_id = ?", account.ID, role.ID).First(&ar).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return e.ErrRoleAlreadyAssigned
	}

	return r.conn.WithContext(ctx).Create(&entities.AccountRole{AccountID: account.ID, RoleID: role.ID}).Error
}

func (r *roleRepository) GetAccountRoles(ctx context.Context, account *entities.Account) ([]*entities.Role, error) {
	var roles []*entities.Role

	if err := r.conn.WithContext(ctx).Joins("join account_roles on account_roles.role_id = roles.id and account_roles.account_id = ?", account.ID).
		Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleRepository) Paginate(ctx context.Context, organizationID uuid.UUID, filters repositories.Filters, order string, limit int, offset int) (roles []*entities.Role, total int64, err error) {
	query := r.conn.WithContext(ctx).Model(&entities.Role{}).
		Where("organization_id = ?", organizationID).
		Or("organization_id is NULL").
		Where(filters.Where).
		Not(filters.Not).
		Preload("Permissions")

	if err = query.Count(&total).Error; err != nil {
		return
	}

	if err = query.Order(order).Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return
	}
	return
}

func (r *roleRepository) Create(ctx context.Context, role *entities.Role) (*entities.Role, error) {
	if err := r.conn.WithContext(ctx).Create(&role).Error; err != nil {
		return nil, err
	}

	return r.FindBy(ctx, map[string]interface{}{"id": role.ID})
}

func (r *roleRepository) Update(ctx context.Context, role *entities.Role) (*entities.Role, error) {
	if err := r.conn.WithContext(ctx).Updates(&role).Error; err != nil {
		return nil, err
	}

	return r.FindBy(ctx, map[string]interface{}{"id": role.ID})
}

func (r *roleRepository) Revoke(ctx context.Context, account *entities.Account, role *entities.Role) error {
	return r.conn.WithContext(ctx).Where("account_id = ? and role_id = ?", account.ID, role.ID).Delete(&entities.AccountRole{}).Error
}

func (r *roleRepository) Delete(ctx context.Context, role *entities.Role) error {
	err := r.conn.WithContext(ctx).Where("role_id = ?", role.ID).First(&entities.AccountRole{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		return e.ErrRoleInUse
	}

	return r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.conn.WithContext(ctx).Where("role_id = ?", role.ID).Delete(&entities.AccountRole{}).Error; err != nil {
			return err
		}

		if err := r.conn.WithContext(ctx).Where("id = ?", role.ID).Delete(&role).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *roleRepository) FindBy(ctx context.Context, params map[string]interface{}) (role *entities.Role, err error) {
	if err = r.conn.WithContext(ctx).Where(params).Preload("Permissions").First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	return
}
