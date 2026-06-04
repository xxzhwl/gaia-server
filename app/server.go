package server

import (
	"os"
	"strings"

	// 开发环境指标样本数据生成器

	"github.com/xxzhwl/gaia/components/asynctask"
	"github.com/xxzhwl/gaia/framework/server"
	"github.com/xxzhwl/gaia/gexit"

	appasynctask "gaia-server/app/asynctask"
	_ "gaia-server/app/jobs" // 注册定时任务 cron service
	"gaia-server/app/metricsdemo"
	"gaia-server/app/router"
	"gaia-server/app/service"
)

func RunServer(port string) {
	if err := service.InitAccount(); err != nil {
		panic("初始化账号服务失败: " + err.Error())
	}
	service.InitAuthService()
	service.InitUserService()
	service.InitRoleService()
	service.InitPermissionService()

	// 注册异步任务和定时任务服务
	appasynctask.RegisterTasks()
	service.GetJobsService().InitJobRunner()

	// 自动启动异步任务调度器
	asynctask.StartScheduler(gexit.GetExitContext(), "GaiaServer")

	// 非生产环境启动指标样本数据生成器（方便 Grafana 大盘预览）
	if env := os.Getenv("GAIA_DISABLE_METRICS_DEMO"); !strings.EqualFold(env, "true") {
		metricsdemo.Start(gexit.GetExitContext())
	}

	s := server.NewAppWithPort(port)
	router.RegisterRouters(s)
	s.Run()
}
