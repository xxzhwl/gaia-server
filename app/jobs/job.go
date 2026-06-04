package jobs

import (
	"errors"
	"fmt"
	"time"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/components/jobs"
	"github.com/xxzhwl/gaia/framework/messageImpl"
)

func init() {
	jobs.RegisterCronService("firstService", &FirstService{})
	jobs.RegisterCronService("secondService", &SecondService{})
	jobs.RegisterCronService("testService", &TestService{})
	jobs.RegisterCronService("logService", &LogService{})
	jobs.RegisterCronService("notifyService", &NotifyService{})
	jobs.RegisterCronService("statsService", &StatsService{})
}

// FirstService 示例服务1
type FirstService struct{}

func (f *FirstService) Start() {
	gaia.Info("CronJobStartFirstService")
}

func (f *FirstService) Stop(arg any) {
	gaia.InfoF("CronJobStopFirstService %v", arg)
}

func (f *FirstService) ReportErr() error {
	return errors.New("CronJobReportErrFirstService")
}

// SecondService 示例服务2
type SecondService struct{}

func (s *SecondService) End() {
	gaia.Info("CronJobEndSecondService")
	robot := messageImpl.NewFeiShuRobot()
	robot.SendText("我发个消息不过分吧")
	time.Sleep(5 * time.Second)
}

func (s *SecondService) Back() {
	for {
		select {
		case <-time.After(time.Second * 5):
			gaia.Info("这是应该后台进程任务")
		}
	}
}

// TestService 测试服务 - 用于测试定时任务功能
type TestService struct{}

// Hello 简单的测试方法
func (t *TestService) Hello() string {
	gaia.Info("TestService.Hello 被调用")
	return "Hello, this is a test job!"
}

// HelloWithArgs 带参数的测试方法
func (t *TestService) HelloWithArgs(name string, count int) string {
	msg := fmt.Sprintf("Hello %s, count: %d", name, count)
	gaia.InfoF("TestService.HelloWithArgs: %s", msg)
	return msg
}

// SimulateWork 模拟耗时工作
func (t *TestService) SimulateWork(duration int) string {
	gaia.InfoF("TestService.SimulateWork 开始工作，预计耗时 %d 秒", duration)
	time.Sleep(time.Duration(duration) * time.Second)
	gaia.InfoF("TestService.SimulateWork 工作完成")
	return fmt.Sprintf("工作完成，耗时 %d 秒", duration)
}

// SimulateError 模拟错误
func (t *TestService) SimulateError() error {
	gaia.Info("TestService.SimulateError 被调用，即将返回错误")
	return errors.New("这是一个模拟的错误")
}

// LogService 日志服务 - 用于测试日志相关任务
type LogService struct{}

// CleanOldLogs 清理旧日志（模拟）
func (l *LogService) CleanOldLogs(days int) string {
	gaia.InfoF("LogService.CleanOldLogs: 清理 %d 天前的日志", days)
	return fmt.Sprintf("已清理 %d 天前的日志", days)
}

// ArchiveLogs 归档日志（模拟）
func (l *LogService) ArchiveLogs() string {
	gaia.Info("LogService.ArchiveLogs: 开始归档日志")
	time.Sleep(2 * time.Second)
	panic("测试panic")
	gaia.Info("LogService.ArchiveLogs: 归档完成")
	return "日志归档完成"
}

// NotifyService 通知服务 - 用于测试通知相关任务
type NotifyService struct{}

// SendDailyReport 发送每日报告（模拟）
func (n *NotifyService) SendDailyReport() string {
	gaia.Info("NotifyService.SendDailyReport: 发送每日报告")
	report := fmt.Sprintf("每日报告 - %s", time.Now().Format("2006-01-02 15:04:05"))

	robot := messageImpl.NewFeiShuRobot()
	robot.SendText(report)

	return report
}

// SendAlert 发送告警（模拟）
func (n *NotifyService) SendAlert(level string, message string) string {
	gaia.InfoF("NotifyService.SendAlert: [%s] %s", level, message)

	robot := messageImpl.NewFeiShuRobot()
	robot.SendText(fmt.Sprintf("【%s告警】%s", level, message))

	return fmt.Sprintf("已发送%s级告警: %s", level, message)
}

// StatsService 统计服务 - 用于测试数据统计任务
type StatsService struct{}

// CalculateDailyStats 计算每日统计（模拟）
func (s *StatsService) CalculateDailyStats() string {
	gaia.Info("StatsService.CalculateDailyStats: 开始计算每日统计")

	time.Sleep(1 * time.Second)

	stats := map[string]any{
		"date":      time.Now().Format("2006-01-02"),
		"users":     100,
		"orders":    500,
		"revenue":   10000.00,
		"timestamp": time.Now().Unix(),
	}

	gaia.InfoF("StatsService.CalculateDailyStats: 统计结果 %+v", stats)
	return fmt.Sprintf("每日统计完成: 用户数=%d, 订单数=%d", 100, 500)
}

// HealthCheck 健康检查
func (s *StatsService) HealthCheck() string {
	gaia.Info("StatsService.HealthCheck: 执行健康检查")

	result := map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"memory":    "正常",
		"cpu":       "正常",
		"database":  "正常",
	}

	gaia.InfoF("StatsService.HealthCheck: 检查结果 %+v", result)
	return "健康检查通过"
}

// BackupData 数据备份（模拟）
func (s *StatsService) BackupData() string {
	gaia.Info("StatsService.BackupData: 开始数据备份")

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	duration := time.Since(startTime).Seconds()
	result := fmt.Sprintf("数据备份完成，耗时 %.2f 秒", duration)
	gaia.InfoF("StatsService.BackupData: %s", result)

	return result
}
