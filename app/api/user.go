package api

import (
	"errors"
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/errwrap"
	"github.com/xxzhwl/gaia/framework/account"
	"github.com/xxzhwl/gaia/framework/server"
)

type UserCtrl struct{}

func NewUserCtrl() *UserCtrl {
	return &UserCtrl{}
}

type GetUserInfoResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Nickname    string   `json:"nickname"`
	Status      string   `json:"status"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
}

type UpdateUserInfoRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateUserInfoResponse struct {
	User account.UserInfo `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" require:"1"`
	NewPassword string `json:"new_password" require:"1"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

type GetUserListRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
}

type GetUserListResponse struct {
	List  []account.UserInfo `json:"list"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
}

func (u *UserCtrl) GetUserInfo() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (GetUserInfoResponse, error) {
		principal, ok := account.GetPrincipal(arg)
		if !ok {
			return GetUserInfoResponse{}, errwrap.Error(account.ErrInvalidToken, errors.New("未授权访问"))
		}
		userInfo, err := service.GetUserService().GetUserInfo(arg.TraceContext, principal.UserID)
		if err != nil {
			return GetUserInfoResponse{}, err
		}
		return GetUserInfoResponse{
			ID:          userInfo.ID,
			Username:    userInfo.Username,
			Email:       userInfo.Email,
			Phone:       userInfo.Phone,
			Nickname:    userInfo.Nickname,
			Status:      userInfo.Status,
			Roles:       userInfo.Roles,
			Permissions: userInfo.Permissions,
			CreatedAt:   userInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	})
}

func (u *UserCtrl) UpdateUserInfo() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (UpdateUserInfoResponse, error) {
		principal, ok := account.GetPrincipal(arg)
		if !ok {
			return UpdateUserInfoResponse{}, errwrap.Error(account.ErrInvalidToken, errors.New("未授权访问"))
		}
		req := UpdateUserInfoRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return UpdateUserInfoResponse{}, errwrap.Error(400, err)
		}
		userInfo, err := service.GetUserService().UpdateProfile(arg.TraceContext, principal.UserID, req.Nickname, req.AvatarURL)
		if err != nil {
			return UpdateUserInfoResponse{}, err
		}
		return UpdateUserInfoResponse{User: *userInfo}, nil
	})
}

func (u *UserCtrl) ChangePassword() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (ChangePasswordResponse, error) {
		principal, ok := account.GetPrincipal(arg)
		if !ok {
			return ChangePasswordResponse{}, errwrap.Error(account.ErrInvalidToken, errors.New("未授权访问"))
		}
		req := ChangePasswordRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return ChangePasswordResponse{}, errwrap.Error(400, err)
		}
		if err := service.GetUserService().ChangePassword(arg.TraceContext, principal.UserID, req.OldPassword, req.NewPassword); err != nil {
			return ChangePasswordResponse{}, err
		}
		return ChangePasswordResponse{Success: true}, nil
	})
}

func (u *UserCtrl) GetUserList() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (GetUserListResponse, error) {
		req := GetUserListRequest{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return GetUserListResponse{}, errwrap.Error(400, err)
		}
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 20
		}
		result, err := service.GetUserService().ListUsers(arg.TraceContext, account.ListUsersRequest{
			Status:   req.Status,
			Keyword:  req.Keyword,
			Page:     req.Page,
			PageSize: req.PageSize,
		})
		if err != nil {
			return GetUserListResponse{}, err
		}
		return GetUserListResponse{
			List:  result.Items,
			Total: result.Total,
			Page:  result.Page,
		}, nil
	})
}

func (u *UserCtrl) CreateUser() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			Username string `json:"username" require:"1"`
			Password string `json:"password" require:"1"`
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		user := &account.User{
			Username: req.Username,
			Nickname: req.Nickname,
			Email:    strPtr(req.Email),
			Phone:    strPtr(req.Phone),
		}
		if err := service.GetUserService().CreateUser(arg.TraceContext, user, req.Password); err != nil {
			return nil, err
		}
		return map[string]any{
			"id":        user.ID,
			"tenant_id": user.TenantID,
			"username":  user.Username,
			"nickname":  user.Nickname,
			"email":     req.Email,
			"phone":     req.Phone,
			"status":    user.Status,
		}, nil
	})
}

func (u *UserCtrl) UpdateUserStatus() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			UserID string `json:"user_id" require:"1"`
			Status string `json:"status" require:"1"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetUserService().UpdateUserStatus(arg.TraceContext, "", req.UserID, req.Status); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func (u *UserCtrl) ResetUserMFA() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		var req struct {
			UserID string `json:"user_id" require:"1"`
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, errwrap.Error(400, err)
		}
		if err := service.GetUserService().ResetUserMFA(arg.TraceContext, "", req.UserID); err != nil {
			return nil, err
		}
		return map[string]any{"success": true}, nil
	})
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
