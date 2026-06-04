package middleware

import (
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func PermissionMiddleware(permissionCode string) app.HandlerFunc {
	return service.GetAccount().Middleware().RequirePermission(permissionCode)
}

func RoleMiddleware(roleCode string) app.HandlerFunc {
	return service.GetAccount().Middleware().RequireRole(roleCode)
}

func AdminMiddleware() app.HandlerFunc {
	return service.GetAccount().Middleware().RequireRole("platform_admin")
}
