package api

import (
	"errors"
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type MFACtrl struct{}

func NewMFACtrl() *MFACtrl {
	return &MFACtrl{}
}

type VerifyTOTPRequest struct {
	Code string `json:"code" require:"1"`
}

type StepUpCompleteRequest struct {
	ChallengeID string `json:"challenge_id" require:"1"`
	Code        string `json:"code" require:"1"`
}

func (c *MFACtrl) SetupTOTP() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		user, err := service.GetAccount().Users().GetByID(arg.TraceContext, p.UserID)
		if err != nil {
			return nil, err
		}
		return service.GetAccount().MFA().SetupTOTP(arg.TraceContext, p.TenantID, p.UserID, user.Email)
	})
}

func (c *MFACtrl) VerifyTOTP() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		req := VerifyTOTPRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		ok, err := service.GetAccount().MFA().VerifyTOTP(arg.TraceContext, account.TOTPVerifyRequest{
			TenantID: p.TenantID,
			UserID:   p.UserID,
			Code:     req.Code,
		})
		return map[string]bool{"ok": ok}, err
	})
}

func (c *MFACtrl) DisableTOTP() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		return map[string]bool{"success": true}, service.GetAccount().MFA().DisableTOTP(arg.TraceContext, p.TenantID, p.UserID, p)
	})
}

func (c *MFACtrl) StartStepUp() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		challengeID, err := service.GetAccount().Auth().RequestStepUp(arg.TraceContext, p)
		if err != nil {
			return nil, err
		}
		return map[string]string{"challenge_id": challengeID}, nil
	})
}

func (c *MFACtrl) CompleteStepUp() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		p, ok := account.GetPrincipal(arg)
		if !ok {
			return nil, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		req := StepUpCompleteRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		return map[string]bool{"success": true}, service.GetAccount().Auth().CompleteStepUp(arg.TraceContext, p, req.ChallengeID, req.Code)
	})
}
