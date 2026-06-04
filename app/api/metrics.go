package api

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/components/asynctask"
	"github.com/xxzhwl/gaia/components/jobs"
	"github.com/xxzhwl/gaia/framework/server"

	"gaia-server/app/service"
)

// MetricsCtrl 指标监控控制器
type MetricsCtrl struct{}

func NewMetricsCtrl() *MetricsCtrl {
	return &MetricsCtrl{}
}

// GetMetricsOverview 返回 jobs + asynctask 综合指标概览
func (c *MetricsCtrl) GetMetricsOverview() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		result := map[string]any{}

		// Jobs 指标快照
		if runner := service.GetJobsService().GetJobRunner(); runner != nil {
			result["jobs"] = runner.SnapshotMetrics()
		}

		// AsyncTask 调度器状态
		result["asynctask"] = asynctask.GetAllSchedulerStatus()

		return result, nil
	})
}

// GetJobsMetrics 返回 jobs 指标快照
func (c *MetricsCtrl) GetJobsMetrics() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		runner := service.GetJobsService().GetJobRunner()
		if runner == nil {
			return map[string]any{"data": nil}, nil
		}
		return map[string]any{"data": runner.SnapshotMetrics()}, nil
	})
}

// GetAsyncTaskMetrics 返回 asynctask 调度器状态
func (c *MetricsCtrl) GetAsyncTaskMetrics() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		return map[string]any{"data": asynctask.GetAllSchedulerStatus()}, nil
	})
}

// GetJobsGlobalMetrics 返回 jobs OTel 全局计数器快照
func (c *MetricsCtrl) GetJobsGlobalMetrics() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		snap := jobs.GetJobsMetrics().Snapshot()
		return map[string]any{"data": snap}, nil
	})
}
