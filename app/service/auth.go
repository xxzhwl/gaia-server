package service

import (
	"context"

	"github.com/xxzhwl/gaia/framework/account"
)

type AuthService struct {
	acct *account.Manager
}

var authService *AuthService

func InitAuthService() {
	authService = &AuthService{acct: GetAccount()}
}

func GetAuthService() *AuthService {
	return authService
}

func (s *AuthService) Login(ctx context.Context, identifier, password, ip, userAgent string) (*account.AuthResult, error) {
	return s.acct.Auth().Login(ctx, account.LoginRequest{
		Identifier: identifier,
		Password:   password,
		IP:         ip,
		UserAgent:  userAgent,
	})
}

func (s *AuthService) Register(ctx context.Context, req account.RegisterRequest) (*account.AuthResult, error) {
	return s.acct.Auth().Register(ctx, req)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, ip, userAgent string) (*account.AuthResult, error) {
	return s.acct.Auth().Refresh(ctx, account.RefreshRequest{
		RefreshToken: refreshToken,
		IP:           ip,
		UserAgent:    userAgent,
	})
}

func (s *AuthService) Logout(ctx context.Context, accessToken, sessionID string) error {
	return s.acct.Auth().Logout(ctx, account.LogoutRequest{
		AccessToken: accessToken,
		SessionID:   sessionID,
	})
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*account.Principal, error) {
	return s.acct.Auth().Validate(ctx, token)
}
