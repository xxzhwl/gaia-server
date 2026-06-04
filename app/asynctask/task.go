package asynctask

import (
	"errors"
	"fmt"
	"time"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/framework/messageImpl"
)

const SystemName = "GaiaServer"

func RegisterTasks() {
	gaia.RegisterProxy(SystemName, "firstService", &FirstService{})
	gaia.RegisterProxy(SystemName, "secondService", &SecondService{})
	gaia.RegisterProxy(SystemName, "testService", &TestService{})
	gaia.RegisterProxy(SystemName, "logService", &LogService{})
	gaia.RegisterProxy(SystemName, "notifyService", &NotifyService{})
	gaia.RegisterProxy(SystemName, "statsService", &StatsService{})
}

type FirstService struct{}

func (f *FirstService) Start() {
	time.Sleep(2 * time.Second)
	gaia.Info("AsyncTaskStartFirstService")
	messageImpl.NewFeiShuRobot().SendRichText(
		"我试一下",
		[]messageImpl.Content{
			{
				Tag:  "text",
				Text: "我试一下",
			},
		},
	)
}

func (f *FirstService) Stop(arg any) {
	gaia.InfoF("AsyncTaskStopFirstService %v", arg)
}

func (f *FirstService) ReportErr() error {
	return errors.New("AsyncTaskReportErrFirstService")
}

type SecondService struct{}

func (f *SecondService) End() {
	gaia.Info("AsyncTaskEndSecondService")
}

// TestService 测试服务
type TestService struct{}

// Hello 简单测试方法
func (t *TestService) Hello() string {
	gaia.Info("TestService.Hello 被调用")
	return "Hello, this is a test async task!"
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
	return errors.New("这是一个模拟的异步任务错误")
}

// SimulatePanic 模拟 panic
func (t *TestService) SimulatePanic() string {
	gaia.Info("TestService.SimulatePanic 被调用，即将 panic")
	panic("这是一个模拟的异步任务 panic")
}

// LogService 日志服务
type LogService struct{}

// CleanOldLogs 清理旧日志
func (l *LogService) CleanOldLogs(days int) string {
	gaia.InfoF("LogService.CleanOldLogs: 清理 %d 天前的日志", days)
	time.Sleep(500 * time.Millisecond)
	return fmt.Sprintf("已清理 %d 天前的日志", days)
}

// ArchiveLogs 归档日志
func (l *LogService) ArchiveLogs() string {
	gaia.Info("LogService.ArchiveLogs: 开始归档日志")
	time.Sleep(1 * time.Second)
	gaia.Info("LogService.ArchiveLogs: 归档完成")
	return "日志归档完成"
}

// NotifyService 通知服务
type NotifyService struct{}

// SendDailyReport 发送每日报告
func (n *NotifyService) SendDailyReport() string {
	gaia.Info("NotifyService.SendDailyReport: 发送每日报告")
	report := fmt.Sprintf("每日报告 - %s", time.Now().Format("2006-01-02 15:04:05"))
	return report
}

// SendAlert 发送告警
func (n *NotifyService) SendAlert(level string, message string) string {
	gaia.InfoF("NotifyService.SendAlert: [%s] %s", level, message)
	return fmt.Sprintf("已发送%s级告警: %s", level, message)
}

// SendFeishuMessage 发送飞书消息
func (n *NotifyService) SendFeishuMessage(title string, content string) string {
	gaia.InfoF("NotifyService.SendFeishuMessage: %s - %s", title, content)
	robot := messageImpl.NewFeiShuRobot()
	robot.SendText(fmt.Sprintf("【%s】%s", title, content))
	return fmt.Sprintf("飞书消息已发送: %s", title)
}

// StatsService 统计服务
type StatsService struct{}

// CalculateDailyStats 计算每日统计
func (s *StatsService) CalculateDailyStats() string {
	gaia.Info("StatsService.CalculateDailyStats: 开始计算每日统计")
	time.Sleep(800 * time.Millisecond)
	return fmt.Sprintf("每日统计完成: 用户数=%d, 订单数=%d", 100, 500)
}

// HealthCheck 健康检查
func (s *StatsService) HealthCheck() string {
	gaia.Info("StatsService.HealthCheck: 执行健康检查")
	return "健康检查通过"
}

// BackupData 数据备份
func (s *StatsService) BackupData() string {
	gaia.Info("StatsService.BackupData: 开始数据备份")
	startTime := time.Now()
	time.Sleep(2 * time.Second)
	duration := time.Since(startTime).Seconds()
	result := fmt.Sprintf("数据备份完成，耗时 %.2f 秒", duration)
	gaia.InfoF("StatsService.BackupData: %s", result)
	return result
}

// GenerateReport 生成报告
func (s *StatsService) GenerateReport(reportType string) string {
	gaia.InfoF("StatsService.GenerateReport: 生成%s报告", reportType)
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("%s报告生成完成", reportType)
}
