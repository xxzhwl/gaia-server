package service

import (
	"context"
	"sync"
	"time"

	"github.com/xxzhwl/gaia/components/asynctask"
)

// AsyncTaskService 异步任务服务
type AsyncTaskService struct {
	mu sync.RWMutex
}

var asyncTaskService *AsyncTaskService
var asyncTaskOnce sync.Once

// GetAsyncTaskService 获取 AsyncTaskService 实例
func GetAsyncTaskService() *AsyncTaskService {
	asyncTaskOnce.Do(func() {
		asyncTaskService = &AsyncTaskService{}
	})
	return asyncTaskService
}

// ListTasks 获取任务列表
func (s *AsyncTaskService) ListTasks(page, pageSize int, taskStatus []string, taskName, systemName string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := asynctask.ListTasksArgs{
		Page:     page,
		PageSize: pageSize,
	}

	if len(taskStatus) > 0 {
		args.TaskStatus = taskStatus
	}

	if taskName != "" {
		args.TaskName = taskName
	}

	if systemName == "" {
		systemName = "default"
	}

	result, err := asynctask.ListTasks(args, systemName, ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"list":  result.List,
		"total": result.Total,
	}, nil
}

// GetTaskDetail 获取任务详情
func (s *AsyncTaskService) GetTaskDetail(taskId int64) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	detail, err := asynctask.GetTaskDetail(taskId, ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"id":                detail.Id,
		"service_name":      detail.ServiceName,
		"method_name":       detail.MethodName,
		"task_name":         detail.TaskName,
		"arg":               detail.Arg,
		"max_retry_time":    detail.MaxRetryTime,
		"system_name":       detail.SystemName,
		"task_status":       detail.TaskStatus,
		"last_run_time":     detail.LastRunTime,
		"last_run_end_time": detail.LastRunEndTime,
		"last_run_duration": detail.LastRunDuration,
		"create_time":       detail.CreateAt,
		"update_time":       detail.UpdateAt,
		"retry_time":        detail.RetryTime,
		"last_result":       detail.LastResult,
		"last_err_msg":      detail.LastErrMsg,
		"log_id":            detail.LogId,
	}, nil
}

// RetryTask 手动重试任务
func (s *AsyncTaskService) RetryTask(taskId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := asynctask.RetryTask(taskId, ctx)
	if err != nil {
		return err
	}

	return nil
}

// CancelTask 取消任务
func (s *AsyncTaskService) CancelTask(taskId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := asynctask.CancelTask(taskId, ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetTaskExecRecords 获取任务执行记录
func (s *AsyncTaskService) GetTaskExecRecords(taskId int64, page, pageSize int) (any, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	records, total, err := asynctask.GetTaskExecRecords(taskId, page, pageSize, ctx)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
