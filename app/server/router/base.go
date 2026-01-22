// Package router 包注释
// @author wanlizhan
// @created 2025-03-17
package router

import (
	api2 "gaia-server/app/server/api"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/xxzhwl/gaia/framework/server"
)

func RegisterRouters(s *server.Server) {
	// 注册健康检查端点

	group := s.Group("/api")
	s.RegisterCommonHandler()

	registerDemo(group)
}

func registerDemo(s *route.RouterGroup) {
	ctrl := api2.NewDemoCtrl()
	g := s.Group("demo")
	g.POST("/demo", ctrl.Demo())
	g.POST("/req", ctrl.Req())
}
