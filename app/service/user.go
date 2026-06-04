package service

import (
	"context"

	"github.com/xxzhwl/gaia/framework/account"
)

type UserService struct {
	acct *account.Manager
}

var userService *UserService

func InitUserService() {
	userService = &UserService{acct: GetAccount()}
}

func GetUserService() *UserService {
	return userService
}

func (s *UserService) GetUserInfo(ctx context.Context, userID string) (*account.UserInfo, error) {
	return s.acct.Users().GetByID(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID, nickname, avatarURL string) (*account.UserInfo, error) {
	return s.acct.Users().UpdateProfile(ctx, account.UpdateProfileRequest{
		UserID:    userID,
		Nickname:  nickname,
		AvatarURL: avatarURL,
	})
}

func (s *UserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	return s.acct.Users().ChangePassword(ctx, account.ChangePasswordRequest{
		UserID:      userID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
}

func (s *UserService) ListUsers(ctx context.Context, req account.ListUsersRequest) (*account.ListUsersResult, error) {
	return s.acct.Admin().ListUsers(ctx, req)
}

func (s *UserService) CreateUser(ctx context.Context, user *account.User, password string) error {
	created, err := s.acct.Admin().CreateUserWithPassword(ctx, account.CreateUserWithPasswordRequest{
		TenantID:  user.TenantID,
		Username:  user.Username,
		Password:  password,
		Nickname:  user.Nickname,
		Email:     stringPtrValue(user.Email),
		Phone:     stringPtrValue(user.Phone),
		Status:    user.Status,
		AvatarURL: user.AvatarURL,
	})
	if err != nil {
		return err
	}
	user.ID = created.ID
	user.TenantID = created.TenantID
	user.Status = created.Status
	return nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, tenantID, userID, status string) error {
	return s.acct.Admin().UpdateUserStatus(ctx, tenantID, userID, status)
}

func (s *UserService) ResetUserMFA(ctx context.Context, tenantID, userID string) error {
	return s.acct.Admin().ResetUserMFA(ctx, tenantID, userID)
}

func stringPtrValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
