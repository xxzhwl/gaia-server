package middleware

import (
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func AuthMiddleware() app.HandlerFunc {
	return service.GetAccount().Middleware().Authenticate()
}

func OptionalAuthMiddleware() app.HandlerFunc {
	return service.GetAccount().Middleware().OptionalAuthenticate()
}
