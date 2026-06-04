package api

import (
	"encoding/json"
	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/framework/server"
)

// HealthCtrl 健康检查控制器
type HealthCtrl struct{}

func NewHealthCtrl() *HealthCtrl {
	return &HealthCtrl{}
}

// CheckHealth 基本健康检查
func (h *HealthCtrl) CheckHealth() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		healthService := service.GetHealthService()
		return healthService.CheckHealth()
	})
}

// CheckHealthDetailed 详细健康检查
func (h *HealthCtrl) CheckHealthDetailed() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		// 从查询参数获取自定义服务列表
		customServicesStr := arg.GetUrlQuery("custom_services")
		var customServices []service.CustomServiceConfig
		if customServicesStr != "" {
			// 解析 JSON 数组
			if err := json.Unmarshal([]byte(customServicesStr), &customServices); err != nil {
				// 解析失败则忽略自定义服务
				customServices = nil
			}
		}

		healthService := service.GetHealthService()
		return healthService.CheckHealthDetailed(customServices)
	})
}

// GetAvailableServices 获取可用服务列表
func (h *HealthCtrl) GetAvailableServices() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		healthService := service.GetHealthService()
		return healthService.GetAvailableServices(), nil
	})
}
