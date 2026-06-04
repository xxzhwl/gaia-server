package api

import (
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type RoleCtrl struct{}

func NewRoleCtrl() *RoleCtrl {
	return &RoleCtrl{}
}

type CreateRoleRequest struct {
	Name        string `json:"name" require:"1"`
	Code        string `json:"code" require:"1"`
	Description string `json:"description"`
}

type AssignUserRoleRequest struct {
	UserID string `json:"user_id" require:"1"`
	RoleID string `json:"role_id" require:"1"`
}

type RemoveUserRoleRequest struct {
	UserID string `json:"user_id" require:"1"`
	RoleID string `json:"role_id" require:"1"`
}

func (r *RoleCtrl) CreateRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := CreateRoleRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		role, err := service.GetRoleService().CreateRole(arg.TraceContext, account.CreateRoleRequest{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"id":          role.ID,
			"name":        role.Name,
			"code":        role.Code,
			"description": role.Description,
			"status":      role.Status,
			"is_system":   role.IsSystem,
		}, nil
	})
}

func (r *RoleCtrl) ListRoles() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		roles, err := service.GetRoleService().ListRoles(arg.TraceContext, "")
		if err != nil {
			return nil, err
		}
		return roles, nil
	})
}

func (r *RoleCtrl) AssignUserRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := AssignUserRoleRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetRoleService().AssignUserRole(arg.TraceContext, req.UserID, req.RoleID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (r *RoleCtrl) UpdateRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			RoleID      string `json:"role_id" require:"1"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetRoleService().UpdateRole(arg.TraceContext, req.RoleID, req.Name, req.Description); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (r *RoleCtrl) DeleteRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			RoleID string `json:"role_id" require:"1"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetRoleService().DeleteRole(arg.TraceContext, req.RoleID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (r *RoleCtrl) GetRolePermissions() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			RoleID string `json:"role_id" require:"1"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		perms, err := service.GetPermissionService().GetRolePermissions(arg.TraceContext, req.RoleID)
		if err != nil {
			return nil, err
		}
		return perms, nil
	})
}

func (r *RoleCtrl) RemoveUserRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := RemoveUserRoleRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetRoleService().RemoveUserRole(arg.TraceContext, req.UserID, req.RoleID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}
