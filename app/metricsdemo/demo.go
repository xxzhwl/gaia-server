// Package metricsdemo 生成示例指标数据，用于验证 Grafana 大盘展示。
// 仅开发/测试环境使用。
package metricsdemo

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/xxzhwl/gaia/components/asynctask"
	"github.com/xxzhwl/gaia/components/jobs"
)

var (
	startOnce sync.Once
	stopCh    chan struct{}
)

var serviceNames = []string{"GaiaServer", "GaiaAdmin"}
var taskThemes = []string{"demo-theme-a", "demo-theme-b", "demo-theme-c"}
var jobNames = []string{"demo-cleanup", "demo-sync", "demo-report"}

// Start 启动后台指标数据生成器，每 30s 生成一轮样本数据。
func Start(ctx context.Context) {
	startOnce.Do(func() {
		registerDemoGauges()
		stopCh = make(chan struct{})
		go generateLoop(ctx)
	})
}

func registerDemoGauges() {
	meter := otel.Meter("gaia-server/metricsdemo")

	// Jobs 异步仪表（loaded_service / loaded_hook）
	jobsMetrics := jobs.GetJobsMetrics()
	meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			obs.ObserveInt64(jobsMetrics.LoadedServiceJobs, int64(3+rand.Intn(5)),
				metric.WithAttributes(attribute.String("replica", "demo-a")))
			obs.ObserveInt64(jobsMetrics.LoadedHookJobs, int64(1+rand.Intn(3)),
				metric.WithAttributes(attribute.String("replica", "demo-b")))
			return nil
		},
		jobsMetrics.LoadedServiceJobs,
		jobsMetrics.LoadedHookJobs,
	)

	// AsyncTask 异步仪表（queue_depth / queue_capacity / worker_total / worker_running）
	asynctaskMetrics := asynctask.GetMetrics()
	meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			for _, theme := range taskThemes {
				for _, rep := range []string{"demo-a", "demo-b"} {
					obs.ObserveInt64(asynctaskMetrics.QueueDepth, int64(rand.Intn(100)),
						metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme), attribute.String("replica", rep)))
					obs.ObserveInt64(asynctaskMetrics.QueueCapacity, int64(200),
						metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme), attribute.String("replica", rep)))
					obs.ObserveInt64(asynctaskMetrics.WorkerTotal, int64(5+rand.Intn(10)),
						metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme), attribute.String("replica", rep)))
					obs.ObserveInt64(asynctaskMetrics.WorkerRunning, int64(3+rand.Intn(5)),
						metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme), attribute.String("replica", rep)))
				}
			}
			return nil
		},
		asynctaskMetrics.QueueDepth,
		asynctaskMetrics.QueueCapacity,
		asynctaskMetrics.WorkerTotal,
		asynctaskMetrics.WorkerRunning,
	)
}

func generateJobsMetrics(ctx context.Context) {
	m := jobs.GetJobsMetrics()

	// ---- 任务执行 ----
	count := 3 + rand.Intn(5)
	for i := 0; i < count; i++ {
		job := jobNames[rand.Intn(len(jobNames))]
		svc := serviceNames[rand.Intn(len(serviceNames))]
		duration := float64(50 + rand.Intn(3000))

		r := rand.Float64()
		var status string
		switch {
		case r < 0.75:
			status = "success"
		case r < 0.90:
			status = "failed"
		case r < 0.97:
			status = "panic"
		default:
			status = "timeout"
		}

		attrs := metric.WithAttributes(
			attribute.String("service_name", svc),
			attribute.String("job_name", job),
			attribute.String("status", status),
			attribute.String("replica", "demo-a"),
		)
		m.JobExecTotal.Add(ctx, 1, attrs)
		m.JobExecDuration.Record(ctx, duration, attrs)

		if status == "panic" {
			m.JobPanicTotal.Add(ctx, 1, metric.WithAttributes(
				attribute.String("service_name", svc),
				attribute.String("job_name", job),
				attribute.String("replica", "demo-a"),
			))
		}
		if status == "timeout" {
			m.JobTimeoutTotal.Add(ctx, 1, metric.WithAttributes(
				attribute.String("service_name", svc),
				attribute.String("job_name", job),
				attribute.String("replica", "demo-a"),
			))
		}
	}

	// ---- Skip ----
	if rand.Float64() < 0.2 {
		m.JobSkipTotal.Add(ctx, 1, metric.WithAttributes(
			attribute.String("service_name", serviceNames[rand.Intn(len(serviceNames))]),
			attribute.String("job_name", jobNames[rand.Intn(len(jobNames))]),
		))
	}

	// ---- DB 错误 ----
	if rand.Float64() < 0.1 {
		m.DBErrorTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("replica", "demo-a")))
	}

	// ---- 告警 ----
	if rand.Float64() < 0.3 {
		m.AlarmFired.Add(ctx, 1)
	}

	// ---- HA 租约 ----
	m.LeaseAcquireTotal.Add(ctx, 1)
	if rand.Float64() < 0.7 {
		m.LeaseAcquiredTotal.Add(ctx, 1)
	} else {
		m.LeaseMissedTotal.Add(ctx, 1)
	}
	if rand.Float64() < 0.8 {
		m.LeaseRenewTotal.Add(ctx, 1)
	}
	if rand.Float64() < 0.05 {
		m.LeaseRenewLost.Add(ctx, 1)
	}

	// ---- updateJobs 扫描 ----
	m.UpdateScanTotal.Add(ctx, 1)
	m.UpdateScanDuration.Record(ctx, float64(10+rand.Intn(200)),
		metric.WithAttributes(attribute.String("replica", "demo-a")))
}

func generateAsyncTaskMetrics(ctx context.Context) {
	m := asynctask.GetMetrics()

	// ---- 任务执行 ----
	count := 5 + rand.Intn(10)
	for i := 0; i < count; i++ {
		theme := taskThemes[rand.Intn(len(taskThemes))]
		duration := float64(100 + rand.Intn(5000))

		r := rand.Float64()
		var status string
		switch {
		case r < 0.7:
			status = "success"
		case r < 0.85:
			status = "failed"
		case r < 0.95:
			status = "retry"
		default:
			status = "panic"
		}

		attrs := metric.WithAttributes(
			attribute.String("status", status),
			asynctask.MetricLabel.Theme.String(theme),
		)

		m.TaskExecTotal.Add(ctx, 1, attrs)
		m.TaskExecDuration.Record(ctx, duration, attrs)

		if status == "panic" {
			m.TaskPanicTotal.Add(ctx, 1, metric.WithAttributes(
				asynctask.MetricLabel.Theme.String(theme),
				asynctask.MetricLabel.Phase.String("exec"),
			))
		}
		if status == "retry" {
			m.TaskRetryTotal.Add(ctx, 1, metric.WithAttributes(
				asynctask.MetricLabel.Theme.String(theme),
			))
		}
	}

	// ---- 扫描 ----
	theme := taskThemes[rand.Intn(len(taskThemes))]
	scanErr := rand.Float64() < 0.08
	m.ScanTotal.Add(ctx, 1, metric.WithAttributes(
		asynctask.MetricLabel.Theme.String(theme),
		attribute.String("status", map[bool]string{true: "failed", false: "success"}[scanErr]),
	))
	m.ScanDuration.Record(ctx, float64(50+rand.Intn(500)),
		metric.WithAttributes(
			attribute.String("status", map[bool]string{true: "failed", false: "success"}[scanErr]),
			asynctask.MetricLabel.Theme.String(theme),
		),
	)

	// ---- 队列 ----
	if rand.Float64() < 0.15 {
		m.QueueDropTotal.Add(ctx, 1, metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}
	m.QueueWaitMillis.Record(ctx, float64(rand.Intn(2000)),
		metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))

	// ---- Worker ----
	if rand.Float64() < 0.3 {
		m.WorkerScaleUp.Add(ctx, int64(1+rand.Intn(3)),
			metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}
	if rand.Float64() < 0.2 {
		m.WorkerScaleDown.Add(ctx, int64(1+rand.Intn(2)),
			metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}

	// ---- 心跳 ----
	if rand.Float64() < 0.05 {
		m.HeartbeatDeadDetected.Add(ctx, 1, metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}

	// ---- DB/告警 ----
	if rand.Float64() < 0.08 {
		m.DBErrorTotal.Add(ctx, 1, metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}
	if rand.Float64() < 0.3 {
		m.AlarmFired.Add(ctx, 1, metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}
	if rand.Float64() < 0.1 {
		m.AlarmSuppressed.Add(ctx, 1, metric.WithAttributes(asynctask.MetricLabel.Theme.String(theme)))
	}
}

func Stop() {
	if stopCh != nil {
		close(stopCh)
	}
}

func generateLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-stopCh:
			return
		case <-ticker.C:
			generateJobsMetrics(ctx)
			generateAsyncTaskMetrics(ctx)
		}
	}
}
