package api

import (
	"errors"
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type SessionCtrl struct{}

func NewSessionCtrl() *SessionCtrl {
	return &SessionCtrl{}
}

func (c *SessionCtrl) List() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		return service.GetAccount().Sessions().List(arg.TraceContext, p.TenantID, p.UserID, p.SessionID)
	})
}

func (c *SessionCtrl) Revoke() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		return map[string]bool{"success": true}, service.GetAccount().Sessions().Revoke(arg.TraceContext, p.TenantID, p.UserID, arg.GetUrlParam("id"))
	})
}

func (c *SessionCtrl) RevokeOthers() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		revoked, err := service.GetAccount().Sessions().RevokeOther(arg.TraceContext, p.TenantID, p.UserID, p.SessionID)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"revoked": revoked}, nil
	})
}
