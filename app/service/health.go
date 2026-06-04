package service

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/components/redis"
)

// HealthService 健康检查服务
type HealthService struct {
	startTime time.Time
}

var healthService *HealthService

func init() {
	healthService = &HealthService{
		startTime: time.Now(),
	}
}

func GetHealthService() *HealthService {
	return healthService
}

// CheckHealth 基本健康检查
func (s *HealthService) CheckHealth() (map[string]any, error) {
	return map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().UnixMilli(),
	}, nil
}

// SystemOverview 系统概览信息
type SystemOverview struct {
	Hostname    string         `json:"hostname"`
	OS          string         `json:"os"`
	Arch        string         `json:"arch"`
	GoVersion   string         `json:"go_version"`
	NumCPU      int            `json:"num_cpu"`
	Uptime      int64          `json:"uptime"`
	UptimeHuman string         `json:"uptime_human"`
	Goroutine   *GoroutineInfo `json:"goroutine"`
	Memory      *MemoryInfo    `json:"memory"`
	GC          *GCInfo        `json:"gc"`
}

// GoroutineInfo 协程信息
type GoroutineInfo struct {
	CurrentCount int `json:"current_count"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	AllocMB      float64 `json:"alloc_mb"`
	TotalAllocMB float64 `json:"total_alloc_mb"`
	SysMB        float64 `json:"sys_mb"`
	HeapAllocMB  float64 `json:"heap_alloc_mb"`
	HeapSysMB    float64 `json:"heap_sys_mb"`
}

// GCInfo GC信息
type GCInfo struct {
	NumGC        uint32  `json:"num_gc"`
	PauseTotalMs float64 `json:"pause_total_ms"`
}

// GetSystemOverview 获取系统概览信息
func (s *HealthService) GetSystemOverview() *SystemOverview {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	hostname, _ := os.Hostname()
	uptime := time.Since(s.startTime)

	return &SystemOverview{
		Hostname:    hostname,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		GoVersion:   runtime.Version(),
		NumCPU:      runtime.NumCPU(),
		Uptime:      int64(uptime.Seconds()),
		UptimeHuman: formatDuration(uptime),
		Goroutine: &GoroutineInfo{
			CurrentCount: runtime.NumGoroutine(),
		},
		Memory: &MemoryInfo{
			AllocMB:      float64(m.Alloc) / 1024 / 1024,
			TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
			SysMB:        float64(m.Sys) / 1024 / 1024,
			HeapAllocMB:  float64(m.HeapAlloc) / 1024 / 1024,
			HeapSysMB:    float64(m.HeapSys) / 1024 / 1024,
		},
		GC: &GCInfo{
			NumGC:        m.NumGC,
			PauseTotalMs: float64(m.PauseTotalNs) / 1e6,
		},
	}
}

// formatDuration 格式化运行时长
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return formatString("%d天%d小时%d分钟", days, hours, minutes)
	}
	if hours > 0 {
		return formatString("%d小时%d分钟", hours, minutes)
	}
	if minutes > 0 {
		return formatString("%d分钟%d秒", minutes, seconds)
	}
	return formatString("%d秒", seconds)
}

func formatString(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// CheckHealthDetailed 详细健康检查
func (s *HealthService) CheckHealthDetailed(customServices []CustomServiceConfig) (map[string]any, error) {
	overview := s.GetSystemOverview()
	services := s.checkServices(customServices)

	return map[string]any{
		"status":    s.getOverallStatus(services),
		"timestamp": time.Now().UnixMilli(),
		"system":    overview,
		"services":  services,
	}, nil
}

// GetAvailableServices 获取可用的服务列表
func (s *HealthService) GetAvailableServices() map[string]any {
	return map[string]any{
		"default": []string{"mysql", "redis"},
		"custom":  []string{},
	}
}

// getSystemInfo 获取系统信息（兼容旧接口）
func (s *HealthService) getSystemInfo() map[string]any {
	overview := s.GetSystemOverview()
	return map[string]any{
		"hostname":      overview.Hostname,
		"os":            overview.OS,
		"arch":          overview.Arch,
		"go_version":    overview.GoVersion,
		"num_cpu":       overview.NumCPU,
		"num_goroutine": overview.Goroutine.CurrentCount,
		"memory_usage": map[string]any{
			"alloc":       overview.Memory.AllocMB,
			"total_alloc": overview.Memory.TotalAllocMB,
			"sys":         overview.Memory.SysMB,
			"num_gc":      overview.GC.NumGC,
		},
		"uptime": overview.Uptime,
	}
}

// CustomServiceConfig 自定义服务配置
type CustomServiceConfig struct {
	Type       string `json:"type" form:"type" query:"type"`                      // mysql 或 redis
	ConfigName string `json:"config_name" form:"config_name" query:"config_name"` // 配置名称，如 Framework.Mysql
}

// checkServices 检查服务状态
func (s *HealthService) checkServices(customServices []CustomServiceConfig) map[string]map[string]any {
	services := make(map[string]map[string]any)

	// 检查默认 MySQL（Framework.Mysql 配置）
	mysqlStatus := s.checkDefaultMySQL()
	if mysqlStatus != nil {
		services["mysql"] = mysqlStatus
	}

	// 检查默认 Redis（Framework.Redis 配置）
	redisStatus := s.checkDefaultRedis()
	if redisStatus != nil {
		services["redis"] = redisStatus
	}

	// 检查自定义服务
	for _, customService := range customServices {
		if status := s.checkCustomService(customService); status != nil {
			serviceKey := customService.ConfigName
			if serviceKey == "" {
				serviceKey = customService.Type + "_custom"
			}
			services[serviceKey] = status
		}
	}

	return services
}

// checkCustomService 检查自定义服务
func (s *HealthService) checkCustomService(service CustomServiceConfig) map[string]any {
	if service.Type == "mysql" {
		return s.checkMySQLByConfig(service.ConfigName)
	}
	if service.Type == "redis" {
		return s.checkRedisByConfig(service.ConfigName)
	}
	return nil
}

// checkMySQLByConfig 根据配置名称检查 MySQL
func (s *HealthService) checkMySQLByConfig(configName string) map[string]any {
	start := time.Now()

	// MySQL 配置直接是 DSN 字符串
	dsn := gaia.GetSafeConfString(configName)
	if dsn == "" {
		return map[string]any{
			"status":      "down",
			"error":       "未找到 MySQL 配置: " + configName,
			"timestamp":   time.Now().UnixMilli(),
			"custom":      true,
			"config_name": configName,
		}
	}

	// 尝试创建连接检查
	mysqlDB, err := gaia.NewMysqlWithSchema(configName)
	if err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"custom":      true,
			"config_name": configName,
		}
	}

	sqlDB, err := mysqlDB.GetGormDb().DB()
	if err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "获取连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"custom":      true,
			"config_name": configName,
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "Ping 失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"custom":      true,
			"config_name": configName,
		}
	}

	return map[string]any{
		"status":      "healthy",
		"latency":     time.Since(start).Milliseconds(),
		"timestamp":   time.Now().UnixMilli(),
		"custom":      true,
		"config_name": configName,
	}
}

// checkRedisByConfig 根据配置名称检查 Redis
func (s *HealthService) checkRedisByConfig(configName string) map[string]any {
	start := time.Now()

	// Redis 配置格式：ConfigName.Address, ConfigName.UserName, ConfigName.Password
	address := gaia.GetSafeConfString(configName + ".Address")
	if address == "" {
		return map[string]any{
			"status":      "down",
			"error":       "未找到 Redis 配置: " + configName + ".Address",
			"timestamp":   time.Now().UnixMilli(),
			"custom":      true,
			"config_name": configName,
		}
	}

	userName := gaia.GetSafeConfString(configName + ".UserName")
	password := gaia.GetSafeConfString(configName + ".Password")

	// 创建 Redis 客户端检查连接
	client := redis.NewClient(address, userName, password)
	cli := client.GetCli()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"custom":      true,
			"config_name": configName,
		}
	}

	return map[string]any{
		"status":      "healthy",
		"latency":     time.Since(start).Milliseconds(),
		"timestamp":   time.Now().UnixMilli(),
		"custom":      true,
		"config_name": configName,
	}
}

// checkDefaultMySQL 检查默认 MySQL 配置
func (s *HealthService) checkDefaultMySQL() map[string]any {
	// MySQL 配置直接是 DSN 字符串，不需要 .Dsn 后缀
	dsn := gaia.GetSafeConfString("Framework.Mysql")
	if dsn == "" {
		// 没有配置 MySQL，返回 nil 表示不显示
		return nil
	}

	return s.checkMySQL()
}

// checkDefaultRedis 检查默认 Redis 配置
func (s *HealthService) checkDefaultRedis() map[string]any {
	address := gaia.GetSafeConfString("Framework.Redis.Address")
	if address == "" {
		// 没有配置 Redis，返回 nil 表示不显示
		return nil
	}

	return s.checkRedis()
}

// checkMySQL 检查 MySQL 状态
func (s *HealthService) checkMySQL() map[string]any {
	start := time.Now()

	mysqlDB, err := gaia.NewFrameworkMysql()
	if err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "获取数据库连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"config_name": "Framework.Mysql",
			"custom":      false,
		}
	}

	sqlDB, err := mysqlDB.GetGormDb().DB()
	if err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "获取数据库连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"config_name": "Framework.Mysql",
			"custom":      false,
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "数据库连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"config_name": "Framework.Mysql",
			"custom":      false,
		}
	}

	latency := time.Since(start).Milliseconds()

	return map[string]any{
		"status":      "healthy",
		"latency":     latency,
		"timestamp":   time.Now().UnixMilli(),
		"config_name": "Framework.Mysql",
		"custom":      false,
	}
}

// checkRedis 检查 Redis 状态
func (s *HealthService) checkRedis() map[string]any {
	redisConfig := gaia.GetSafeConfString("Framework.Redis.Address")
	if redisConfig == "" {
		return nil
	}

	start := time.Now()

	userName := gaia.GetSafeConfString("Framework.Redis.UserName")
	password := gaia.GetSafeConfString("Framework.Redis.Password")

	// 创建 Redis 客户端检查连接
	client := redis.NewClient(redisConfig, userName, password)
	cli := client.GetCli()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		return map[string]any{
			"status":      "down",
			"error":       "连接失败: " + err.Error(),
			"timestamp":   time.Now().UnixMilli(),
			"config_name": "Framework.Redis",
			"custom":      false,
		}
	}

	return map[string]any{
		"status":      "healthy",
		"latency":     time.Since(start).Milliseconds(),
		"timestamp":   time.Now().UnixMilli(),
		"config_name": "Framework.Redis",
		"custom":      false,
	}
}

// getOverallStatus 获取整体状态
func (s *HealthService) getOverallStatus(services map[string]map[string]any) string {
	for _, service := range services {
		if service["status"] == "down" {
			return "unhealthy"
		}
	}

	for _, service := range services {
		if service["status"] == "warning" {
			return "warning"
		}
	}

	return "healthy"
}
