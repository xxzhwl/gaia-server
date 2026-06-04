package api

import (
	"errors"
	"gaia-server/app/service"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type AuthCtrl struct{}

func NewAuthCtrl() *AuthCtrl {
	return &AuthCtrl{}
}

type LoginRequest struct {
	Identifier     string `json:"identifier" require:"1"`
	IdentifierType string `json:"identifier_type"`
	Password       string `json:"password" require:"1"`
}

type LoginResponse struct {
	User                 account.UserInfo `json:"user"`
	AccessToken          string           `json:"access_token"`
	RefreshToken         string           `json:"refresh_token"`
	ExpiresAt            string           `json:"expires_at"`
	TokenType            string           `json:"token_type"`
	MFARequired          bool             `json:"mfa_required"`
	MFAChallenge         string           `json:"mfa_challenge,omitempty"`
	PhoneBindingRequired bool             `json:"phone_binding_required"`
}

type RegisterRequest struct {
	Username string `json:"username" require:"1"`
	Password string `json:"password" require:"1"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type RegisterResponse LoginResponse

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" require:"1"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

type LogoutRequest struct{}

type LogoutResponse struct {
	Success bool `json:"success"`
}

type ForgotPasswordStartRequest struct {
	TenantID   string `json:"tenant_id"`
	Identifier string `json:"identifier" require:"1"`
	Channel    string `json:"channel"`
}

type ForgotPasswordCompleteRequest struct {
	TenantID    string `json:"tenant_id"`
	ChallengeID string `json:"challenge_id" require:"1"`
	Code        string `json:"code" require:"1"`
	NewPassword string `json:"new_password" require:"1"`
}

type CompleteMFARequest struct {
	ChallengeID string `json:"challenge_id" require:"1"`
	Code        string `json:"code" require:"1"`
}

func (a *AuthCtrl) Login() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (LoginResponse, error) {
		req := LoginRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return LoginResponse{}, errwrap.Error(400, err)
		}
		authSvc := service.GetAuthService()
		result, err := authSvc.Login(arg.TraceContext, req.Identifier, req.Password, arg.C().ClientIP(), string(arg.C().UserAgent()))
		if err != nil {
			return LoginResponse{}, err
		}
		return LoginResponse{
			User:                 *result.User,
			AccessToken:          result.AccessToken,
			RefreshToken:         result.RefreshToken,
			ExpiresAt:            formatTime(result.ExpiresAt),
			TokenType:            result.TokenType,
			MFARequired:          result.MFARequired,
			MFAChallenge:         result.MFAChallenge,
			PhoneBindingRequired: result.PhoneBindingRequired,
		}, nil
	})
}

func (a *AuthCtrl) Register() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (RegisterResponse, error) {
		req := RegisterRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return RegisterResponse{}, errwrap.Error(400, err)
		}
		authSvc := service.GetAuthService()
		result, err := authSvc.Register(arg.TraceContext, account.RegisterRequest{
			Username:  req.Username,
			Password:  req.Password,
			Nickname:  req.Nickname,
			Email:     req.Email,
			Phone:     req.Phone,
			DeviceID:  string(arg.C().GetHeader("X-Device-ID")),
			IP:        arg.C().ClientIP(),
			UserAgent: string(arg.C().UserAgent()),
		})
		if err != nil {
			return RegisterResponse{}, err
		}
		return RegisterResponse{
			User:         *result.User,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    formatTime(result.ExpiresAt),
			TokenType:    result.TokenType,
		}, nil
	})
}

func (a *AuthCtrl) RefreshToken() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (RefreshTokenResponse, error) {
		req := RefreshTokenRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return RefreshTokenResponse{}, errwrap.Error(400, err)
		}
		authSvc := service.GetAuthService()
		result, err := authSvc.RefreshToken(arg.TraceContext, req.RefreshToken, arg.C().ClientIP(), string(arg.C().UserAgent()))
		if err != nil {
			return RefreshTokenResponse{}, err
		}
		return RefreshTokenResponse{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    formatTime(result.ExpiresAt),
			TokenType:    result.TokenType,
		}, nil
	})
}

func (a *AuthCtrl) Logout() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (LogoutResponse, error) {
		principal, ok := account.GetPrincipal(arg)
		if !ok {
			return LogoutResponse{}, errwrap.Error(account.ErrInvalidToken, errors.New("未认证"))
		}
		authSvc := service.GetAuthService()
		if err := authSvc.Logout(arg.TraceContext, "", principal.SessionID); err != nil {
			return LogoutResponse{}, err
		}
		return LogoutResponse{Success: true}, nil
	})
}

func (a *AuthCtrl) CompleteMFA() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (LoginResponse, error) {
		req := CompleteMFARequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return LoginResponse{}, errwrap.Error(400, err)
		}
		result, err := service.GetAccount().Auth().CompleteMFA(arg.TraceContext, req.ChallengeID, req.Code, account.CompleteMFARequest{
			DeviceID:  string(arg.C().GetHeader("X-Device-ID")),
			IP:        arg.C().ClientIP(),
			UserAgent: string(arg.C().UserAgent()),
		})
		if err != nil {
			return LoginResponse{}, err
		}
		return LoginResponse{
			User:         *result.User,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresAt:    formatTime(result.ExpiresAt),
			TokenType:    result.TokenType,
		}, nil
	})
}

func (a *AuthCtrl) ForgotPasswordStart() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := ForgotPasswordStartRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		return service.GetAccount().Users().StartResetPassword(arg.TraceContext, account.StartResetPasswordRequest{
			TenantID:   req.TenantID,
			Identifier: req.Identifier,
			Channel:    req.Channel,
		})
	})
}

func (a *AuthCtrl) ForgotPasswordComplete() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		req := ForgotPasswordCompleteRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		return map[string]bool{"success": true}, service.GetAccount().Users().CompleteResetPassword(arg.TraceContext, account.CompleteResetPasswordRequest{
			TenantID:    req.TenantID,
			ChallengeID: req.ChallengeID,
			Code:        req.Code,
			NewPassword: req.NewPassword,
		})
	})
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
