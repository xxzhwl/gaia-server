package api

import (
	"strconv"

	"gaia-server/app/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/framework/server"
)

// JobsCtrl 定时任务控制器
type JobsCtrl struct{}

func NewJobsCtrl() *JobsCtrl {
	return &JobsCtrl{}
}

// ListJobs 获取任务列表
func (c *JobsCtrl) ListJobs() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		pageStr := arg.GetUrlQuery("page")
		pageSizeStr := arg.GetUrlQuery("page_size")
		jobType := arg.GetUrlQuery("job_type")
		jobName := arg.GetUrlQuery("job_name")
		enabledStr := arg.GetUrlQuery("enabled")

		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		pageSize := 20
		if pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
				pageSize = ps
			}
		}

		var enabled *bool
		if enabledStr != "" {
			e := enabledStr == "1" || enabledStr == "true"
			enabled = &e
		}

		svc := service.GetJobsService()
		return svc.ListJobs(page, pageSize, jobType, jobName, enabled)
	})
}

// GetJobDetail 获取任务详情
func (c *JobsCtrl) GetJobDetail() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		jobIdStr := arg.GetUrlParam("id")
		jobId, err := strconv.ParseInt(jobIdStr, 10, 64)
		if err != nil {
			return nil, err
		}

		svc := service.GetJobsService()
		return svc.GetJobDetail(jobId)
	})
}

// CreateJob 创建任务
func (c *JobsCtrl) CreateJob() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		req := service.CreateJobReq{
			Timeout: 300,
			Enabled: true,
		}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}

		svc := service.GetJobsService()
		jobId, err := svc.CreateJob(req)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"job_id":  jobId,
			"message": "创建成功",
		}, nil
	})
}

// UpdateJob 更新任务
func (c *JobsCtrl) UpdateJob() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		req := service.UpdateJobReq{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}

		svc := service.GetJobsService()
		affected, err := svc.UpdateJob(req)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"affected": affected,
			"message":  "更新成功",
		}, nil
	})
}

// DeleteJob 删除任务
func (c *JobsCtrl) DeleteJob() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		jobIdStr := arg.GetUrlParam("id")
		jobId, err := strconv.ParseInt(jobIdStr, 10, 64)
		if err != nil {
			return nil, err
		}

		svc := service.GetJobsService()
		err = svc.DeleteJob(jobId)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"message": "删除成功",
		}, nil
	})
}

// ToggleJob 启用/禁用任务
func (c *JobsCtrl) ToggleJob() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		type ToggleJobReq struct {
			JobId  int64 `json:"job_id" require:"1"`
			Enable bool  `json:"enable" require:"1"`
		}

		req := ToggleJobReq{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}

		svc := service.GetJobsService()
		err := svc.ToggleJob(req.JobId, req.Enable)
		if err != nil {
			return nil, err
		}

		status := "禁用"
		if req.Enable {
			status = "启用"
		}

		return map[string]any{
			"message": status + "成功",
		}, nil
	})
}

// ExecuteJobImmediately 立即执行任务
func (c *JobsCtrl) ExecuteJobImmediately() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		type ExecuteReq struct {
			JobId int64 `json:"job_id" require:"1"`
		}

		req := ExecuteReq{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}

		svc := service.GetJobsService()
		err := svc.ExecuteJobImmediately(req.JobId)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"message": "任务已触发执行",
		}, nil
	})
}

// GetJobRecords 获取任务执行记录
func (c *JobsCtrl) GetJobRecords() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		jobIdStr := arg.GetUrlQuery("job_id")
		pageStr := arg.GetUrlQuery("page")
		pageSizeStr := arg.GetUrlQuery("page_size")

		var jobId int64
		if jobIdStr != "" {
			if id, err := strconv.ParseInt(jobIdStr, 10, 64); err == nil {
				jobId = id
			}
		}

		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		pageSize := 20
		if pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
				pageSize = ps
			}
		}

		svc := service.GetJobsService()
		records, total, err := svc.GetJobRecords(jobId, page, pageSize)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"data": records,
			"sum":  total,
		}, nil
	})
}

// GetRunningJobs 获取正在运行的任务
func (c *JobsCtrl) GetRunningJobs() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		svc := service.GetJobsService()
		runningJobs, err := svc.GetRunningJobs()
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"data": runningJobs,
		}, nil
	})
}

// GetJobTypes 获取任务类型列表
func (c *JobsCtrl) GetJobTypes() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		return map[string]any{
			"data": []map[string]string{
				{"value": "cron_job", "label": "服务方法调用"},
				{"value": "cron_hook", "label": "HTTP Hook"},
			},
		}, nil
	})
}

// GetJobStats 获取任务统计
func (c *JobsCtrl) GetJobStats() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		svc := service.GetJobsService()
		stats, err := svc.GetJobStats()
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"data": stats,
		}, nil
	})
}

// StopExecutor 暂停 Jobs 执行器
func (c *JobsCtrl) StopExecutor() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		svc := service.GetJobsService()
		err := svc.StopExecutor()
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"message": "执行器已暂停",
		}, nil
	})
}

// StartExecutor 启动 Jobs 执行器
func (c *JobsCtrl) StartExecutor() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		svc := service.GetJobsService()
		err := svc.StartExecutor()
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"message": "执行器已启动",
		}, nil
	})
}
