package api

import (
	"strconv"

	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type PermissionCtrl struct{}

func NewPermissionCtrl() *PermissionCtrl {
	return &PermissionCtrl{}
}

func (p *PermissionCtrl) Create() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			ResourceType string `json:"resource_type" require:"1"`
			Action       string `json:"action" require:"1"`
			Code         string `json:"code" require:"1"`
			Description  string `json:"description"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		perm := &account.Permission{
			ResourceType: req.ResourceType,
			Action:       req.Action,
			Code:         req.Code,
			Description:  req.Description,
		}
		if err := service.GetPermissionService().Create(arg.TraceContext, perm); err != nil {
			return nil, err
		}
		return map[string]any{
			"id":            perm.ID,
			"resource_type": perm.ResourceType,
			"action":        perm.Action,
			"code":          perm.Code,
			"description":   perm.Description,
			"status":        perm.Status,
		}, nil
	})
}

func (p *PermissionCtrl) Update() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			PermissionID string `json:"permission_id" require:"1"`
			Description  string `json:"description"`
			Status       string `json:"status"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetPermissionService().Update(arg.TraceContext, "", req.PermissionID, req.Description, req.Status); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (p *PermissionCtrl) Delete() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			PermissionID string `json:"permission_id" require:"1"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetPermissionService().Delete(arg.TraceContext, "", req.PermissionID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (p *PermissionCtrl) List() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := account.PermissionListRequest{
			Page:     1,
			PageSize: 200,
		}
		if p := arg.GetUrlQuery("page"); p != "" {
			if v, err := strconv.Atoi(p); err == nil {
				req.Page = v
			}
		}
		if s := arg.GetUrlQuery("page_size"); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				req.PageSize = v
			}
		}
		req.ResourceType = arg.GetUrlQuery("resource_type")
		req.Action = arg.GetUrlQuery("action")
		return service.GetPermissionService().List(arg.TraceContext, req)
	})
}

func (p *PermissionCtrl) AssignToRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			RoleID       string `json:"role_id" require:"1"`
			PermissionID string `json:"permission_id" require:"1"`
			TenantID     string `json:"tenant_id"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetPermissionService().AssignToRole(arg.TraceContext, req.TenantID, req.RoleID, req.PermissionID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (p *PermissionCtrl) RemoveFromRole() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			RoleID       string `json:"role_id" require:"1"`
			PermissionID string `json:"permission_id" require:"1"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetPermissionService().RemoveFromRole(arg.TraceContext, req.RoleID, req.PermissionID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}
