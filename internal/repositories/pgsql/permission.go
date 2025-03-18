package pgsql

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"context"
	"errors"
	"gorm.io/gorm"
)

type permissionRepository struct {
	conn *gorm.DB
}

func NewPermissionRepository(conn *gorm.DB) *permissionRepository {
	return &permissionRepository{
		conn: conn,
	}
}

func (r *permissionRepository) FindBy(ctx context.Context, params map[string]interface{}) (permission *entities.Permission, err error) {
	if err = r.conn.WithContext(ctx).Where(params).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrEntityNotFound
		}

		return
	}

	return
}

func (r *permissionRepository) Create(ctx context.Context, permission *entities.Permission) (*entities.Permission, error) {
	err := r.conn.WithContext(ctx).Create(permission).Error
	if err != nil {
		return nil, err
	}

	return r.FindBy(ctx, map[string]interface{}{"id": permission.ID})
}

func (r *permissionRepository) All(ctx context.Context) (permissions []*entities.Permission, err error) {
	err = r.conn.WithContext(ctx).Find(&permissions).Error

	return
}

func (r *permissionRepository) Update(ctx context.Context, permission *entities.Permission) (*entities.Permission, error) {
	err := r.conn.WithContext(ctx).Updates(&permission).Error
	if err != nil {
		return nil, err
	}

	return r.FindBy(ctx, map[string]interface{}{"id": permission.ID})
}

func (r *permissionRepository) Delete(ctx context.Context, permission *entities.Permission) error {
	return r.conn.WithContext(ctx).Where("id = ?", permission.ID).Delete(&permission).Error
}

func (r *permissionRepository) GetAccountPermissions(ctx context.Context, account *entities.Account) ([]*entities.Permission, error) {
	var permissions []*entities.Permission

	if err := r.conn.WithContext(ctx).Joins("join account_permissions on account_permissions.permission_id = permissions.id and account_permissions.account_id = ?", account.ID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *permissionRepository) Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) (permissions []*entities.Permission, total int64, err error) {
	query := r.conn.WithContext(ctx).Model(&entities.Permission{}).Where(filters)

	if err = query.Count(&total).Error; err != nil {
		return
	}

	if err = query.Order(order).Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		return
	}
	return
}

//func (r *permissionRepository) RevokeAccountPermissions(ctx context.Context, account *entities.Account, permissions ...*entities.Permission) error {
//	return r.conn.WithContext(ctx).Model(&account).Association("Permissions").Delete(&permissions)
//}

func (r *permissionRepository) RevokeRolePermissions(ctx context.Context, role *entities.Role, permissions ...*entities.Permission) error {
	return r.conn.WithContext(ctx).Model(&role).Association("Permissions").Delete(&permissions)
}

//func (r *permissionRepository) AssignAccountPermissions(ctx context.Context, account *entities.Account, permissions []*entities.Permission) error {
//	return r.conn.WithContext(ctx).Model(&account).Omit("UpdatedAt").Association("Permissions").Append(&permissions)
//}

func (r *permissionRepository) AssignRolePermissions(ctx context.Context, role *entities.Role, permissions []*entities.Permission) error {
	return r.conn.WithContext(ctx).Model(&role).Association("Permissions").Append(&permissions)
}
