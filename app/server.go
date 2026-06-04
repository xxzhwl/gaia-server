package server

import (
	"os"
	"strings"
	"time"

	// 开发环境指标样本数据生成器

	"github.com/xxzhwl/gaia"
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
	start := time.Now()

	gaia.InfoF("开始初始化 Account 模块...")
	if err := service.InitAccount(); err != nil {
		panic("初始化账号服务失败: " + err.Error())
	}
	gaia.InfoF("Account 模块初始化完成，耗时: %v", time.Since(start))

	service.InitAuthService()
	service.InitUserService()
	service.InitRoleService()
	service.InitPermissionService()

	// 注册异步任务和定时任务服务
	appasynctask.RegisterTasks()

	jobsStart := time.Now()
	service.GetJobsService().InitJobRunner()
	gaia.InfoF("Jobs 模块初始化完成，耗时: %v", time.Since(jobsStart))

	// 自动启动异步任务调度器（先执行表迁移）
	asyncStart := time.Now()
	asyncScheduler := asynctask.NewScheduler("GaiaServer")
	if err := asyncScheduler.Bootstrap(gexit.GetExitContext()); err != nil {
		panic("asynctask bootstrap failed: " + err.Error())
	}
	gaia.InfoF("AsyncTask 模块初始化完成，耗时: %v", time.Since(asyncStart))

	asynctask.StartScheduler(gexit.GetExitContext(), "GaiaServer")

	// 非生产环境启动指标样本数据生成器（方便 Grafana 大盘预览）
	if env := os.Getenv("GAIA_DISABLE_METRICS_DEMO"); !strings.EqualFold(env, "true") {
		metricsdemo.Start(gexit.GetExitContext())
	}

	gaia.InfoF("所有服务初始化完成，总耗时: %v", time.Since(start))

	s := server.NewAppWithPort(port)
	router.RegisterRouters(s)
	s.Run()
}
