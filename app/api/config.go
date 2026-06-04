package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xxzhwl/gaia"
	"github.com/xxzhwl/gaia/framework/remoteConfig"
	"github.com/xxzhwl/gaia/framework/server"
)

type ConfigCtrl struct{}

func NewConfigCtrl() *ConfigCtrl {
	return &ConfigCtrl{}
}

// GetLocalConfig 返回本地配置文件内容（支持 .json / .yaml / .yml）
func (c *ConfigCtrl) GetLocalConfig() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		configDir := filepath.Dir(gaia.DefaultLocalConfigFile)

		// 按优先级尝试不同格式的配置文件
		type candidate struct {
			path   string
			format string
		}
		candidates := []candidate{
			{filepath.Join(configDir, "config.yaml"), "yaml"},
			{filepath.Join(configDir, "config.yml"), "yaml"},
			{gaia.DefaultLocalConfigFile, "json"},
		}

		var data []byte
		var fileFormat string
		var err error
		for _, c := range candidates {
			data, err = gaia.ReadFileAll(c.path)
			if err == nil {
				fileFormat = c.format
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("读取本地配置文件失败（已尝试 json/yaml/yml）: %w", err)
		}

		if fileFormat == "yaml" {
			return map[string]any{
				"__content": string(data),
				"__format":  "yaml",
			}, nil
		}

		var parsed any
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("解析本地配置文件失败: %w", err)
		}
		return parsed, nil
	})
}

// GetRemoteConfig 返回远程配置中心的内容
func (c *ConfigCtrl) GetRemoteConfig() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		// 优先从远程配置中心实时获取
		if center := remoteConfig.ActiveCenter(); center != nil {
			// 优先尝试全量导出（支持 DumpConfig 的 Provider 如 Nacos）
			if dumper, ok := center.(remoteConfig.ConfigDumper); ok {
				full, err := dumper.DumpConfig()
				if err == nil && full != nil {
					return full, nil
				}
			}

			// 回退：尝试获取顶层 key（配置中心一般按 key 维度存储）
			path := gaia.GetSafeConfString("RemoteConfig.Path")
			if path != "" {
				val, existed, err := center.GetConfig(path)
				if err != nil {
					return nil, fmt.Errorf("从配置中心读取失败: %w", err)
				}
				if existed {
					return val, nil
				}
			}

			return map[string]any{
				"note": "配置中心不支持全量列举，请指定具体 key 查询",
			}, nil
		}

		// 回退到本地缓存文件
		data, err := gaia.ReadFileAll(gaia.DefaultRemoteConfigFile)
		if err != nil {
			return map[string]any{
				"note": "远程配置中心未启用，且本地缓存文件不存在",
			}, nil
		}
		var parsed any
		json.Unmarshal(data, &parsed)
		return parsed, nil
	})
}

// GetEnvironmentConfig 返回环境变量中与 gaia 相关的配置
func (c *ConfigCtrl) GetEnvironmentConfig() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		envVars := map[string]string{}
		for _, e := range os.Environ() {
			if pair := strings.SplitN(e, "=", 2); len(pair) == 2 {
				key := pair[0]
				// 只返回与框架相关的环境变量
				if strings.HasPrefix(key, "GAIA_") ||
					strings.HasPrefix(key, "Deploy") ||
					key == "SystemEnName" ||
					key == "REAL_RUN_ENV" {
					envVars[key] = pair[1]
				}
			}
		}
		return envVars, nil
	})
}

// GetConfigByKey 查询指定 key 的配置（多层来源）
func (c *ConfigCtrl) GetConfigByKey() app.HandlerFunc {
	return server.MakeHandler(func(arg server.Request) (any, error) {
		key := arg.GetUrlQuery("key")
		if key == "" {
			return nil, fmt.Errorf("key 不能为空")
		}

		// 1. 环境变量
		envVal, envExists := os.LookupEnv(key)

		// 2. 本地文件
		localVal, localExists, _ := gaia.GetConfFromLocalFile(key)

		// 3. 远程配置中心
		var remoteVal any
		var remoteExists bool
		if center := remoteConfig.ActiveCenter(); center != nil {
			v, e, _ := center.GetConfig(key)
			remoteVal = v
			remoteExists = e
		}

		// 4. 最终生效值
		finalVal, _ := gaia.GetConf(key)

		return map[string]any{
			"key":    key,
			"env":    cond(envExists, envVal, nil),
			"local":  cond(localExists, localVal, nil),
			"remote": cond(remoteExists, remoteVal, nil),
			"final":  finalVal,
		}, nil
	})
}

func cond(exists bool, val any, fallback any) any {
	if exists {
		return val
	}
	return fallback
}
