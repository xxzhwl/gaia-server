package api

import (
	"errors"
	"gaia-server/app/middleware"
	"gaia-server/app/service"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type AdminCtrl struct{}

func NewAdminCtrl() *AdminCtrl {
	return &AdminCtrl{}
}

func (c *AdminCtrl) RequirePermission(code string) app.HandlerFunc {
	return middleware.PermissionMiddleware(code)
}

func (c *AdminCtrl) ListUserSessions() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未授权"))
		}
		tenantID := arg.GetUrlQuery("tenant_id")
		if tenantID == "" {
			tenantID = p.TenantID
		}
		return service.GetAccount().Admin().ListUserSessions(arg.TraceContext, tenantID, arg.GetUrlParam("id"))
	})
}

func (c *AdminCtrl) GetUserPermissions() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		if _, ok := account.GetPrincipal(arg); !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未授权"))
		}
		userID := arg.GetUrlParam("id")
		perms, err := service.GetAccount().Authorizer().GetEffectivePermissions(arg.TraceContext, userID)
		if err != nil {
			return nil, err
		}
		return map[string]any{"user_id": userID, "permissions": perms}, nil
	})
}

func (c *AdminCtrl) RevokeSession() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未授权"))
		}
		tenantID := arg.GetUrlQuery("tenant_id")
		if tenantID == "" {
			tenantID = p.TenantID
		}
		return nil, service.GetAccount().Admin().RevokeSession(arg.TraceContext, tenantID, arg.GetUrlParam("id"))
	})
}

func (c *AdminCtrl) QueryAudit() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未授权"))
		}
		req, err := bindAuditQuery(arg, p.TenantID)
		if err != nil {
			return nil, err
		}
		return service.GetAccount().Audit().Query(arg.TraceContext, req)
	})
}

func (c *AdminCtrl) QueryArchivedAudit() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未授权"))
		}
		req, err := bindAuditQuery(arg, p.TenantID)
		if err != nil {
			return nil, err
		}
		return service.GetAccount().Audit().QueryArchived(arg.TraceContext, req)
	})
}

func (c *AdminCtrl) ListAuditEvents() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未授权"))
		}
		tenantID := arg.GetUrlQuery("tenant_id")
		if tenantID == "" {
			tenantID = p.TenantID
		}
		now := time.Now()
		return service.GetAccount().Audit().ListEvents(arg.TraceContext, tenantID, now.AddDate(0, 0, -30), now)
	})
}

func (c *AdminCtrl) GetAuditLog() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		return service.GetAccount().Audit().GetByID(arg.TraceContext, arg.GetUrlParam("id"))
	})
}

func (c *AdminCtrl) RestoreAuditLog() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		return map[string]bool{"success": true}, service.GetAccount().Audit().RestoreFromArchive(arg.TraceContext, arg.GetUrlParam("id"))
	})
}

func bindAuditQuery(arg server.Request, defaultTenantID string) (account.AuditQueryRequest, error) {
	req := account.AuditQueryRequest{
		TenantID: arg.GetUrlQuery("tenant_id"),
		UserID:   arg.GetUrlQuery("user_id"),
		Event:    arg.GetUrlQuery("event"),
		Status:   arg.GetUrlQuery("status"),
		Keyword:  arg.GetUrlQuery("keyword"),
		Page:     parseIntQuery(arg, "page", 1),
		PageSize: parseIntQuery(arg, "page_size", 50),
	}
	if req.TenantID == "" {
		req.TenantID = defaultTenantID
	}
	if start := arg.GetUrlQuery("start_at"); start != "" {
		t, err := parseAuditTime(start)
		if err != nil {
			return req, errwrap.Error(400, err)
		}
		req.StartAt = &t
	}
	if end := arg.GetUrlQuery("end_at"); end != "" {
		t, err := parseAuditTime(end)
		if err != nil {
			return req, errwrap.Error(400, err)
		}
		req.EndAt = &t
	}
	return req, nil
}

func parseIntQuery(arg server.Request, key string, fallback int) int {
	raw := arg.GetUrlQuery(key)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}

func parseAuditTime(raw string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", raw, time.Local); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02", raw, time.Local)
}
