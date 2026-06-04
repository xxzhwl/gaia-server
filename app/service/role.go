package service

import (
	"context"

	"github.com/xxzhwl/gaia/framework/account"
)

type RoleService struct {
	acct *account.Manager
}

var roleService *RoleService

func InitRoleService() {
	roleService = &RoleService{acct: GetAccount()}
}

func GetRoleService() *RoleService {
	return roleService
}

func (s *RoleService) CreateRole(ctx context.Context, req account.CreateRoleRequest) (*account.Role, error) {
	return s.acct.Admin().CreateRole(ctx, req)
}

func (s *RoleService) UpdateRole(ctx context.Context, roleID, name, description string) error {
	return s.acct.Admin().UpdateRole(ctx, roleID, name, description)
}

func (s *RoleService) DeleteRole(ctx context.Context, roleID string) error {
	return s.acct.Admin().DeleteRole(ctx, roleID)
}

func (s *RoleService) ListRoles(ctx context.Context, tenantID string) ([]account.Role, error) {
	return s.acct.Admin().ListRoles(ctx, tenantID)
}

func (s *RoleService) AssignUserRole(ctx context.Context, userID, roleID string) error {
	return s.acct.Admin().AssignUserRole(ctx, userID, roleID)
}

func (s *RoleService) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	return s.acct.Admin().RemoveUserRole(ctx, userID, roleID)
}
