package api

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia/components/asynctask"
	"github.com/xxzhwl/gaia/framework/server"

	"gaia-server/app/service"
)

// AsyncTaskCtrl 异步任务控制器
type AsyncTaskCtrl struct{}

func NewAsyncTaskCtrl() *AsyncTaskCtrl {
	return &AsyncTaskCtrl{}
}

// SchedulerStatus SSE 流式返回调度器状态
func (c *AsyncTaskCtrl) SchedulerStatus() app.HandlerFunc {
	return func(ctx context.Context, hertzCtx *app.RequestContext) {
		hertzCtx.SetContentType("text/event-stream")
		hertzCtx.Response.Header.Set("Cache-Control", "no-cache")
		hertzCtx.Response.Header.Set("Connection", "keep-alive")
		hertzCtx.Response.Header.Set("Access-Control-Allow-Origin", "*")

		statusMap := asynctask.GetAllSchedulerStatus()

		if len(statusMap) == 0 {
			data := map[string]any{
				"Scans":          0,
				"PushTasks":      0,
				"PullTasks":      0,
				"ExecTasks":      0,
				"ExecSuccess":    0,
				"ExecFails":      0,
				"RunningWorkers": 0,
				"AllWorkers":     0,
			}
			jsonData, _ := json.Marshal(data)
			hertzCtx.WriteString("event: timestamp\n")
			hertzCtx.WriteString("data: " + string(jsonData) + "\n\n")
			hertzCtx.Flush()
			return
		}

		var totalStatus = map[string]int64{
			"Scans":          0,
			"PushTasks":      0,
			"PullTasks":      0,
			"ExecTasks":      0,
			"ExecSuccess":    0,
			"ExecFails":      0,
			"RunningWorkers": 0,
			"AllWorkers":     0,
		}

		for _, status := range statusMap {
			totalStatus["Scans"] += int64(status.Scans)
			totalStatus["PushTasks"] += int64(status.PushTasks)
			totalStatus["PullTasks"] += int64(status.PullTasks)
			totalStatus["ExecTasks"] += int64(status.ExecTasks)
			totalStatus["ExecSuccess"] += int64(status.ExecSuccess)
			totalStatus["ExecFails"] += int64(status.ExecFails)
			totalStatus["RunningWorkers"] += int64(status.RunningWorkers)
			totalStatus["AllWorkers"] += int64(status.AllWorkers)
		}

		data := map[string]any{
			"Scans":          totalStatus["Scans"],
			"PushTasks":      totalStatus["PushTasks"],
			"PullTasks":      totalStatus["PullTasks"],
			"ExecTasks":      totalStatus["ExecTasks"],
			"ExecSuccess":    totalStatus["ExecSuccess"],
			"ExecFails":      totalStatus["ExecFails"],
			"RunningWorkers": totalStatus["RunningWorkers"],
			"AllWorkers":     totalStatus["AllWorkers"],
		}

		jsonData, _ := json.Marshal(data)
		hertzCtx.WriteString("event: timestamp\n")
		hertzCtx.WriteString("data: " + string(jsonData) + "\n\n")
		hertzCtx.Flush()
	}
}

// GetSchedulerInfoPoll HTTP 轮询获取调度器信息
func (c *AsyncTaskCtrl) GetSchedulerInfoPoll() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		infoList := asynctask.GetAllSchedulerInfo()

		var totalStatus = map[string]int64{
			"Scans":          0,
			"PushTasks":      0,
			"PullTasks":      0,
			"ExecTasks":      0,
			"ExecSuccess":    0,
			"ExecFails":      0,
			"RunningWorkers": 0,
			"AllWorkers":     0,
		}

		type schedulerStatusResp struct {
			Scans          int64 `json:"scans"`
			PushTasks      int64 `json:"push_tasks"`
			PullTasks      int64 `json:"pull_tasks"`
			ExecTasks      int64 `json:"exec_tasks"`
			ExecSuccess    int64 `json:"exec_success"`
			ExecFails      int64 `json:"exec_fails"`
			RunningWorkers int32 `json:"running_workers"`
			AllWorkers     int32 `json:"all_workers"`
		}

		type schedulerInfoResp struct {
			Theme   string              `json:"theme"`
			Running bool                `json:"running"`
			Status  schedulerStatusResp `json:"status"`
		}

		infos := make([]schedulerInfoResp, 0, len(infoList))
		for _, info := range infoList {
			totalStatus["Scans"] += int64(info.Status.Scans)
			totalStatus["PushTasks"] += int64(info.Status.PushTasks)
			totalStatus["PullTasks"] += int64(info.Status.PullTasks)
			totalStatus["ExecTasks"] += int64(info.Status.ExecTasks)
			totalStatus["ExecSuccess"] += int64(info.Status.ExecSuccess)
			totalStatus["ExecFails"] += int64(info.Status.ExecFails)
			totalStatus["RunningWorkers"] += int64(info.Status.RunningWorkers)
			totalStatus["AllWorkers"] += int64(info.Status.AllWorkers)

			infos = append(infos, schedulerInfoResp{
				Theme:   info.Theme,
				Running: info.Running,
				Status: schedulerStatusResp{
					Scans:          info.Status.Scans,
					PushTasks:      info.Status.PushTasks,
					PullTasks:      info.Status.PullTasks,
					ExecTasks:      info.Status.ExecTasks,
					ExecSuccess:    info.Status.ExecSuccess,
					ExecFails:      info.Status.ExecFails,
					RunningWorkers: info.Status.RunningWorkers,
					AllWorkers:     info.Status.AllWorkers,
				},
			})
		}

		return map[string]any{
			"scans":           totalStatus["Scans"],
			"push_tasks":      totalStatus["PushTasks"],
			"pull_tasks":      totalStatus["PullTasks"],
			"exec_tasks":      totalStatus["ExecTasks"],
			"exec_success":    totalStatus["ExecSuccess"],
			"exec_fails":      totalStatus["ExecFails"],
			"running_workers": totalStatus["RunningWorkers"],
			"all_workers":     totalStatus["AllWorkers"],
			"schedulers":      infos,
		}, nil
	})
}

// ListTasks 获取任务列表
func (c *AsyncTaskCtrl) ListTasks() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		pageStr := arg.GetUrlQuery("page")
		pageSizeStr := arg.GetUrlQuery("page_size")
		taskName := arg.GetUrlQuery("task_name")
		taskStatus := arg.GetUrlQueryArray("task_status")

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

		svc := service.GetAsyncTaskService()
		return svc.ListTasks(page, pageSize, taskStatus, taskName, "GaiaServer")
	})
}

// GetTaskDetail 获取任务详情
func (c *AsyncTaskCtrl) GetTaskDetail() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		taskIdStr := arg.GetUrlParam("id")
		taskId, err := strconv.ParseInt(taskIdStr, 10, 64)
		if err != nil {
			return nil, err
		}

		svc := service.GetAsyncTaskService()
		return svc.GetTaskDetail(taskId)
	})
}

// RetryTask 手动重试任务
func (c *AsyncTaskCtrl) RetryTask() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		type RetryTaskReq struct {
			TaskId int64 `json:"task_id" require:"1"`
		}

		req := RetryTaskReq{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}

		svc := service.GetAsyncTaskService()
		err := svc.RetryTask(req.TaskId)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"message": "重试任务已提交",
		}, nil
	})
}

// CancelTask 取消任务
func (c *AsyncTaskCtrl) CancelTask() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		type CancelTaskReq struct {
			TaskId int64 `json:"task_id" require:"1"`
		}

		req := CancelTaskReq{}
		if err := arg.BindJsonWithChecker(&req); err != nil {
			return nil, err
		}

		svc := service.GetAsyncTaskService()
		err := svc.CancelTask(req.TaskId)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"message": "任务已取消",
		}, nil
	})
}

// GetTaskExecRecords 获取任务执行记录
func (c *AsyncTaskCtrl) GetTaskExecRecords() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		taskIdStr := arg.GetUrlQuery("task_id")
		pageStr := arg.GetUrlQuery("page")
		pageSizeStr := arg.GetUrlQuery("page_size")

		var taskId int64 = 0
		if taskIdStr != "" {
			id, err := strconv.ParseInt(taskIdStr, 10, 64)
			if err != nil {
				return nil, err
			}
			taskId = id
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

		svc := service.GetAsyncTaskService()
		records, total, err := svc.GetTaskExecRecords(taskId, page, pageSize)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"data": records,
			"sum":  total,
		}, nil
	})
}

// GetAllSchedulerStatus 获取所有调度器状态
func (c *AsyncTaskCtrl) GetAllSchedulerStatus() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		statusMap := asynctask.GetAllSchedulerStatus()

		result := make(map[string]any)
		for name, status := range statusMap {
			result[name] = map[string]any{
				"Scans":          status.Scans,
				"PushTasks":      status.PushTasks,
				"PullTasks":      status.PullTasks,
				"ExecTasks":      status.ExecTasks,
				"ExecSuccess":    status.ExecSuccess,
				"ExecFails":      status.ExecFails,
				"RunningWorkers": status.RunningWorkers,
				"AllWorkers":     status.AllWorkers,
			}
		}

		return map[string]any{
			"data": result,
		}, nil
	})
}

func (c *AsyncTaskCtrl) StopScheduler() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		type StopReq struct {
			Theme string `json:"theme"`
		}
		req := StopReq{}
		if err := arg.BindJson(&req); err != nil {
			return nil, err
		}

		if req.Theme == "" {
			return nil, errors.New("theme不能为空")
		}

		scheduler := asynctask.GetScheduler(req.Theme)
		if scheduler == nil {
			return nil, errors.New("未查询到对应的调度器")
		}
		scheduler.Stop()
		return map[string]any{"message": "调度器已停止", "theme": req.Theme}, nil
	})
}

func (c *AsyncTaskCtrl) StartScheduler() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (map[string]any, error) {
		type StartReq struct {
			Theme string `json:"theme"`
		}
		req := StartReq{}
		if err := arg.BindJson(&req); err != nil {
			return nil, err
		}

		if req.Theme == "" {
			return nil, errors.New("theme不能为空")
		}

		scheduler := asynctask.GetScheduler(req.Theme)
		if scheduler == nil {
			return nil, errors.New("未查询到对应的调度器")
		}
		scheduler.Resume()
		return map[string]any{"message": "调度器已启动", "theme": req.Theme}, nil
	})
}
