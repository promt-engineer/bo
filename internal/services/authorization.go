package services

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/internal/repositories"
	"backoffice/internal/transport/http/requests"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const (
	defaultRoleType = "default"
)

var (
	ErrCanNotAssignRole             = errors.New("can not assign role")
	ErrCanNotRevokeAdminPermissions = errors.New("cant revoke admin permissions")
)

type AuthorizationService struct {
	accountService       *AccountService
	roleRepository       repositories.RoleRepository
	permissionRepository repositories.PermissionRepository
}

func NewAuthorizationService(accountService *AccountService, roleRepository repositories.RoleRepository, permissionRepository repositories.PermissionRepository) *AuthorizationService {
	return &AuthorizationService{
		accountService:       accountService,
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
	}
}

func (s *AuthorizationService) LoadAccountRoles(ctx context.Context, account *entities.Account) error {
	roles, err := s.roleRepository.GetAccountRoles(ctx, account)
	if err != nil {
		return err
	}

	account.Roles = roles

	return nil
}

func (s *AuthorizationService) LoadAccountPermissions(ctx context.Context, account *entities.Account) error {
	permissions, err := s.permissionRepository.GetAccountPermissions(ctx, account)
	if err != nil {
		return err
	}

	account.Permissions = permissions

	return nil
}

func (s *AuthorizationService) PaginatePermissions(ctx context.Context, filters map[string]interface{}, order string, limit int, offset int) ([]*entities.Permission, int64, error) {
	return s.permissionRepository.Paginate(ctx, filters, order, limit, offset)
}

func (s *AuthorizationService) CreatePermission(ctx context.Context, name, description, subject, endpoint, action string) (*entities.Permission, error) {
	_, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"name": name})
	if err == nil {
		return nil, e.ErrEntityAlreadyExist
	}

	if err != nil && !errors.Is(err, e.ErrEntityNotFound) {
		return nil, err
	}

	return s.permissionRepository.Create(ctx, &entities.Permission{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Subject:     subject,
		Endpoint:    endpoint,
		Action:      action,
	})
}

func (s *AuthorizationService) GetPermission(ctx context.Context, permissionID uuid.UUID) (*entities.Permission, error) {
	return s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
}

func (s *AuthorizationService) UpdatePermission(ctx context.Context, permissionID uuid.UUID, req *requests.UpsertPermissionRequest) (*entities.Permission, error) {
	permission, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
	if err != nil {
		return nil, err
	}

	return s.permissionRepository.Update(ctx, &entities.Permission{
		ID:          permission.ID,
		Name:        req.Name,
		Description: req.Description,
		Subject:     req.Subject,
		Endpoint:    req.Endpoint,
		Action:      req.Action,
	})
}

func (s *AuthorizationService) DeletePermission(ctx context.Context, permissionID uuid.UUID) error {
	permission, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
	if err != nil {
		return err
	}

	return s.permissionRepository.Delete(ctx, permission)
}

func (s *AuthorizationService) PaginateRoles(ctx context.Context, organizationID uuid.UUID, where map[string]interface{}, order string, limit int, offset int) ([]*entities.Role, int64, error) {
	if where == nil {
		where = make(map[string]interface{})
	}

	filters := repositories.Filters{Where: where, Not: map[string]interface{}{"type": entities.RootRoleTypeName}}

	// by default user can see roles without organization (where organization is null)
	return s.roleRepository.Paginate(ctx, organizationID, filters, order, limit, offset)
}

func (s *AuthorizationService) CreateRole(ctx context.Context, session *entities.Session, name, description, t string) (*entities.Role, error) {
	if !session.Account.IsRoot() {
		t = defaultRoleType
	}

	return s.roleRepository.Create(ctx, &entities.Role{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		Type:           t,
		OrganizationID: session.OrganizationID,
	})
}

func (s *AuthorizationService) UpdateRole(ctx context.Context, session *entities.Session, roleID uuid.UUID, name, description, t string) (*entities.Role, error) {
	if !session.Account.IsRoot() {
		t = defaultRoleType
	}

	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID})
	if err != nil {
		return nil, err
	}

	return s.roleRepository.Update(ctx, &entities.Role{
		ID:             role.ID,
		Name:           name,
		Description:    description,
		Type:           t,
		OrganizationID: session.OrganizationID,
	})
}

func (s *AuthorizationService) GetRole(ctx context.Context, roleID uuid.UUID) (*entities.Role, error) {
	return s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID})
}

func (s *AuthorizationService) CanAssignRole(ctx context.Context, roleID string) error {
	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID})
	if err != nil {
		return err
	}

	if role.Type == entities.RootRoleTypeName {
		return fmt.Errorf("%v: %v", ErrCanNotAssignRole, role.Type)
	}

	return nil
}

func (s *AuthorizationService) AssignRole(ctx context.Context, accountID, roleID string) error {
	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return err
	}

	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID})
	if err != nil {
		return err
	}

	return s.roleRepository.Assign(ctx, account, role)
}

func (s *AuthorizationService) RevokeRole(ctx context.Context, accountID, roleID string) error {
	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
	if err != nil {
		return err
	}

	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID})
	if err != nil {
		return err
	}

	return s.roleRepository.Revoke(ctx, account, role)
}

//func (s *AuthorizationService) AssignAccountPermissions(ctx context.Context, accountID string, permissionIDs ...string) error {
//	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
//	if err != nil {
//		return err
//	}
//
//	var permissions []*entities.Permission
//
//	for _, permissionID := range permissionIDs {
//		permission, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
//		if err != nil {
//			return err
//		}
//
//		permissions = append(permissions, permission)
//	}
//
//	return s.permissionRepository.AssignAccountPermissions(ctx, account, permissions)
//}

func (s *AuthorizationService) AssignRolePermissions(ctx context.Context, organizationID uuid.UUID, roleID string, permissionIDs ...string) error {
	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID, "organization_id": organizationID})
	if err != nil {
		return err
	}

	var permissions []*entities.Permission

	for _, permissionID := range permissionIDs {
		permission, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
		if err != nil {
			return err
		}

		permissions = append(permissions, permission)
	}

	return s.permissionRepository.AssignRolePermissions(ctx, role, permissions)
}

//func (s *AuthorizationService) RevokeAccountPermissions(ctx context.Context, accountID string, permissionIDs ...string) error {
//	account, err := s.accountService.FindBy(ctx, map[string]interface{}{"id": accountID})
//	if err != nil {
//		return err
//	}
//
//	var permissions []*entities.Permission
//
//	for _, permissionID := range permissionIDs {
//		permission, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
//		if err != nil {
//			return err
//		}
//
//		permissions = append(permissions, permission)
//	}
//
//	return s.permissionRepository.RevokeAccountPermissions(ctx, account, permissions...)
//}

func (s *AuthorizationService) RevokeRolePermissions(ctx context.Context, session *entities.Session, roleID string, permissionIDs ...string) error {
	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID, "organization_id": session.OrganizationID})
	if err != nil {
		return err
	}

	var permissions []*entities.Permission

	if role.Type == entities.AdminRoleTypeName && !session.Account.IsRoot() {
		return ErrCanNotRevokeAdminPermissions
	}

	for _, permissionID := range permissionIDs {
		permission, err := s.permissionRepository.FindBy(ctx, map[string]interface{}{"id": permissionID})
		if err != nil {
			return err
		}

		permissions = append(permissions, permission)
	}

	return s.permissionRepository.RevokeRolePermissions(ctx, role, permissions...)
}

func (s *AuthorizationService) DeleteRole(ctx context.Context, organizationID uuid.UUID, roleID string) error {
	role, err := s.roleRepository.FindBy(ctx, map[string]interface{}{"id": roleID, "organization_id": organizationID})
	if err != nil {
		return err
	}

	return s.roleRepository.Delete(ctx, role)
}
