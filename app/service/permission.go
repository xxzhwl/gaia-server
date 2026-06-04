package service

import (
	"context"

	"github.com/xxzhwl/gaia/framework/account"
)

type PermissionService struct {
	acct *account.Manager
}

var permissionService *PermissionService

func InitPermissionService() {
	permissionService = &PermissionService{acct: GetAccount()}
}

func GetPermissionService() *PermissionService {
	return permissionService
}

func (s *PermissionService) List(ctx context.Context, req account.PermissionListRequest) (*account.PermissionListResult, error) {
	return s.acct.Permissions().List(ctx, req)
}

func (s *PermissionService) Create(ctx context.Context, perm *account.Permission) error {
	return s.acct.Admin().CreatePermission(ctx, perm)
}

func (s *PermissionService) Update(ctx context.Context, tenantID, permissionID, description, status string) error {
	return s.acct.Admin().UpdatePermission(ctx, tenantID, permissionID, description, status)
}

func (s *PermissionService) Delete(ctx context.Context, tenantID, permissionID string) error {
	return s.acct.Admin().DeletePermission(ctx, tenantID, permissionID)
}

func (s *PermissionService) GetRolePermissions(ctx context.Context, roleID string) ([]account.Permission, error) {
	return s.acct.Admin().GetRolePermissions(ctx, roleID)
}

func (s *PermissionService) AssignToRole(ctx context.Context, tenantID, roleID, permissionID string) error {
	return s.acct.Permissions().AssignToRole(ctx, tenantID, roleID, permissionID)
}

func (s *PermissionService) RemoveFromRole(ctx context.Context, roleID, permissionID string) error {
	return s.acct.Permissions().RemoveFromRole(ctx, roleID, permissionID)
}
