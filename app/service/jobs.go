package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/xxzhwl/gaia/components/jobs"
)

// JobsService 定时任务服务
type JobsService struct {
	mu        sync.RWMutex
	jobRunner *jobs.RunJob
}

var jobsService *JobsService
var jobsOnce sync.Once

// GetJobsService 获取 JobsService 实例
func GetJobsService() *JobsService {
	jobsOnce.Do(func() {
		jobsService = &JobsService{}
	})
	return jobsService
}

// InitJobRunner 初始化 JobRunner
func (s *JobsService) InitJobRunner() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.jobRunner == nil {
		s.jobRunner = jobs.NewSecondJobs()
		go s.jobRunner.Run()
	}
}

// GetJobRunner 获取 JobRunner 实例
func (s *JobsService) GetJobRunner() *jobs.RunJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobRunner
}

// StopExecutor 暂停 Jobs 执行器
func (s *JobsService) StopExecutor() error {
	runner := s.GetJobRunner()
	if runner == nil {
		return errors.New("Jobs 调度器未初始化")
	}
	runner.Stop()
	return nil
}

// StartExecutor 启动 Jobs 执行器
func (s *JobsService) StartExecutor() error {
	runner := s.GetJobRunner()
	if runner == nil {
		return errors.New("Jobs 调度器未初始化")
	}
	runner.Resume()
	return nil
}

// ListJobs 获取任务列表
func (s *JobsService) ListJobs(page, pageSize int, jobType, jobName string, enabled *bool) (map[string]any, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return nil, errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := jobs.ListJobsArgs{
		Page:     page,
		PageSize: pageSize,
	}

	if jobType != "" {
		args.JobType = jobType
	}

	if jobName != "" {
		args.JobName = jobName
	}

	if enabled != nil {
		args.Enabled = enabled
	}

	result, err := runner.ListJobs(args, ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"list":  result.List,
		"total": result.Total,
	}, nil
}

// GetJobDetail 获取任务详情
func (s *JobsService) GetJobDetail(jobId int64) (map[string]any, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return nil, errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	detail, err := runner.GetJobDetail(jobId, ctx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"id":             detail.Id,
		"job_name":       detail.JobName,
		"job_type":       detail.JobType,
		"cron_expr":      detail.CronExpr,
		"hook_url":       detail.HookUrl,
		"service_name":   detail.ServiceName,
		"service_method": detail.ServiceMethod,
		"args":           string(detail.Args),
		"timeout":        detail.Timeout,
		"enabled":        detail.Enabled,
		"run_status":     detail.RunStatus,
		"last_run_time":  detail.LastRunTime,
		"create_time":    detail.CreateTime,
		"update_time":    detail.UpdateTime,
	}, nil
}

// CreateJobReq 创建任务请求
type CreateJobReq struct {
	JobName       string `json:"job_name"`
	JobType       string `json:"job_type"`
	CronExpr      string `json:"cron_expr"`
	ServiceName   string `json:"service_name"`
	ServiceMethod string `json:"service_method"`
	HookUrl       string `json:"hook_url"`
	Args          string `json:"args"`
	Timeout       int    `json:"timeout"`
	Enabled       bool   `json:"enabled"`
}

// CreateJob 创建任务
func (s *JobsService) CreateJob(req CreateJobReq) (int64, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return 0, errors.New("Jobs 调度器未初始化")
	}

	if req.JobType != jobs.CronJob && req.JobType != jobs.CronHook {
		return 0, errors.New("无效的任务类型")
	}

	var args []byte
	if req.Args != "" {
		args = []byte(req.Args)
	} else {
		args = []byte("{}")
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 300
	}

	addArgs := jobs.AddJobArgs{
		JobBase: jobs.JobBase{
			JobName:       req.JobName,
			JobType:       req.JobType,
			CronExpr:      req.CronExpr,
			ServiceName:   req.ServiceName,
			ServiceMethod: req.ServiceMethod,
			HookUrl:       req.HookUrl,
			Args:          args,
			Timeout:       timeout,
		},
		Enable: req.Enabled,
	}

	jobId, err := runner.AddJob(addArgs)
	if err != nil {
		return 0, err
	}

	return jobId, nil
}

// UpdateJobReq 更新任务请求
type UpdateJobReq struct {
	JobId         int64  `json:"job_id"`
	JobName       string `json:"job_name"`
	CronExpr      string `json:"cron_expr"`
	ServiceName   string `json:"service_name"`
	ServiceMethod string `json:"service_method"`
	HookUrl       string `json:"hook_url"`
	Args          string `json:"args"`
	Timeout       int    `json:"timeout"`
	Enabled       *bool  `json:"enabled"`
}

// UpdateJob 更新任务
func (s *JobsService) UpdateJob(req UpdateJobReq) (int64, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return 0, errors.New("Jobs 调度器未初始化")
	}

	var args []byte
	if req.Args != "" {
		args = []byte(req.Args)
	} else {
		args = []byte("{}")
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 300
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	updateArgs := jobs.UpdateJobArgs{
		JobId: req.JobId,
		JobBase: jobs.JobBase{
			JobName:       req.JobName,
			CronExpr:      req.CronExpr,
			ServiceName:   req.ServiceName,
			ServiceMethod: req.ServiceMethod,
			HookUrl:       req.HookUrl,
			Args:          args,
			Timeout:       timeout,
		},
		Enable: enabled,
	}

	affected, err := runner.UpdateJob(updateArgs)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// DeleteJob 删除任务
func (s *JobsService) DeleteJob(jobId int64) error {
	runner := s.GetJobRunner()
	if runner == nil {
		return errors.New("Jobs 调度器未初始化")
	}

	err := runner.RemoveJob(jobId)
	if err != nil {
		return err
	}

	return nil
}

// ToggleJob 启用/禁用任务
func (s *JobsService) ToggleJob(jobId int64, enable bool) error {
	runner := s.GetJobRunner()
	if runner == nil {
		return errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := runner.ToggleJob(jobId, enable, ctx)
	if err != nil {
		return err
	}

	return nil
}

// ExecuteJobImmediately 立即执行任务
func (s *JobsService) ExecuteJobImmediately(jobId int64) error {
	runner := s.GetJobRunner()
	if runner == nil {
		return errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := runner.ExecuteJobImmediately(jobId, ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetJobRecords 获取任务执行记录
func (s *JobsService) GetJobRecords(jobId int64, page, pageSize int) (any, int64, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return nil, 0, errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	records, total, err := runner.GetJobRecords(jobId, page, pageSize, ctx)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetRunningJobs 获取正在运行的任务
func (s *JobsService) GetRunningJobs() (any, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return nil, errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	runningJobs, err := runner.GetRunningJobs(ctx)
	if err != nil {
		return nil, err
	}

	return runningJobs, nil
}

// JobStats 任务统计
type JobStats struct {
	TotalJobs     int `json:"total_jobs"`
	EnabledJobs   int `json:"enabled_jobs"`
	DisabledJobs  int `json:"disabled_jobs"`
	RunningJobs   int `json:"running_jobs"`
	CronJobCount  int `json:"cron_job_count"`
	CronHookCount int `json:"cron_hook_count"`
}

// GetJobStats 获取任务统计
func (s *JobsService) GetJobStats() (*JobStats, error) {
	runner := s.GetJobRunner()
	if runner == nil {
		return nil, errors.New("Jobs 调度器未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := runner.ListJobs(jobs.ListJobsArgs{
		Page:     1,
		PageSize: 10000,
	}, ctx)
	if err != nil {
		return nil, err
	}

	runningJobs, err := runner.GetRunningJobs(ctx)
	if err != nil {
		return nil, err
	}

	stats := &JobStats{
		RunningJobs: len(runningJobs),
	}

	for _, job := range result.List {
		stats.TotalJobs++

		if job.Enabled {
			stats.EnabledJobs++
		} else {
			stats.DisabledJobs++
		}

		if job.JobType == jobs.CronJob {
			stats.CronJobCount++
		} else if job.JobType == jobs.CronHook {
			stats.CronHookCount++
		}
	}

	return stats, nil
}
