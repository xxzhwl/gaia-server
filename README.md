# Gaia Server

基于 Gaia SDK 构建的高性能 Go HTTP 服务，提供完整的后台服务框架，支持 HTTP API、异步任务、定时任务等多种服务模式。

## 特性

- 🚀 **高性能 HTTP 服务**：基于 CloudWeGo Hertz 框架，提供高性能的 HTTP API 服务
- 🔧 **模块化设计**：采用 Gaia SDK 封装，提供配置管理、缓存、数据库、日志等核心模块
- 📦 **多服务管理**：支持 HTTP Server、异步任务、定时任务等多种服务类型，可通过 Supervisord 统一管理
- ⚙️ **灵活配置**：支持本地 JSON 配置和远程配置中心，配置结构清晰易于维护
- 🔄 **内置通用功能**：包含 JWT 认证、跨域支持、健康检查、性能监控等常用功能
- 🐳 **容器化部署**：提供完整的 Docker 部署方案，支持开发和生产环境

## 技术栈

- **编程语言**: Go 1.24.4
- **Web 框架**: CloudWeGo Hertz
- **ORM**: GORM + GORM Gen
- **数据库**: MySQL, PostgreSQL, ClickHouse
- **缓存**: Redis
- **搜索**: Elasticsearch
- **对象存储**: 腾讯云 COS
- **进程管理**: Supervisord
- **容器化**: Docker
- **监控追踪**: OpenTelemetry + Jaeger
- **消息通知**: 飞书机器人

## 项目结构

```
gaia-server/
├── app/server/                    # 应用核心代码
│   ├── api/                       # API 控制器
│   │   └── demo.go                # 示例 API
│   ├── asynctask/                 # 异步任务处理
│   │   └── task.go                # 异步任务注册
│   ├── jobs/                      # 定时任务
│   │   └── job.go                 # 定时任务定义
│   ├── repo/                      # 数据访问层
│   │   ├── dao/                   # GORM Gen 生成的 DAO
│   │   └── model/                 # 数据模型
│   ├── router/                    # 路由定义
│   │   └── base.go                # 路由注册
│   └── server.go                  # 服务启动入口
├── bin/                           # 可执行文件入口
│   └── service.go                 # 主服务入口，支持多服务模式
├── cmd/                           # 服务启动脚本
│   └── run_http_server.sh         # HTTP 服务启动脚本
├── configs/                       # 配置文件
│   ├── common/                    # 通用配置（操作/查询）
│   ├── local/                     # 本地开发配置
│   │   └── config.json            # 主配置文件
│   ├── remote/                    # 远程配置缓存
│   └── readme.md                  # 配置说明文档
├── deploy/                        # 部署相关文件
│   ├── Dockerfile                 # Docker 构建文件
│   ├── supervisord.conf           # Supervisord 配置
│   └── deploy.sh                  # 部署脚本
├── var/logs/                      # 日志目录
├── go.mod                         # Go 模块定义
├── go.sum                         # 依赖校验
└── .dockerignore                  # Docker 忽略文件
```

## 快速开始

### 环境要求

- Go 1.24.4 或更高版本
- MySQL 5.7+ / PostgreSQL 12+ / Redis 6+
- Docker (可选，用于容器化部署)

### 本地运行

1. **克隆项目并安装依赖**
   ```bash
   git clone <项目地址>
   cd gaia-server
   go mod download
   ```

2. **配置数据库和缓存**
   复制配置文件模板并根据实际情况修改：
   ```bash
   cp configs/local/config.json.example configs/local/config.json
   ```
   编辑 `configs/local/config.json`，配置数据库连接、Redis 地址等。

3. **启动 HTTP 服务**
   ```bash
   # 直接运行
   go run ./bin -Service=Server
   
   # 或使用启动脚本
   ./cmd/run_http_server.sh
   ```

4. **验证服务**
   服务默认运行在 `http://localhost:8008`
   ```bash
   curl http://localhost:8008/health
   ```

### API 示例

项目包含一个示例 API，展示 Gaia SDK 的基本用法：

**Demo API**
- `POST /api/demo/demo` - 处理 JSON 请求并返回数据
- `POST /api/demo/req`   - 演示 HTTP 客户端调用

示例请求：
```bash
curl -X POST http://localhost:8008/api/demo/demo \
  -H "Content-Type: application/json" \
  -d '{"name": "Gaia User"}'
```

响应：
```json
{"Name": "Gaia User"}
```

## 配置说明

详细的配置说明请参考 [configs/readme.md](configs/readme.md)。

### 主要配置项

```json
{
  "SystemEnName": "GaiaServer",          // 系统英文名称
  "SystemCnName": "GaiaServer",          // 系统中文名称
  "Framework": {                         // 框架核心配置
    "Mysql": "DSN连接字符串",            // MySQL 数据库
    "Redis": {                           // Redis 配置
      "Address": "host:port",
      "Password": "密码"
    },
    "ES": {                              // Elasticsearch 配置
      "Address": "http://host:port",
      "UserName": "用户名",
      "Password": "密码"
    },
    "Cos": {                             // 腾讯云 COS 配置
      "appId": "应用ID",
      "secretId": "密钥ID",
      "secretKey": "密钥"
    }
  },
  "Server": {                            // HTTP 服务器配置
    "Port": "8008",                      // 服务端口
    "Cors": {                            // 跨域配置
      "Enable": false,
      "AllowOrigins": ["http://localhost:3002"]
    },
    "Logger": {                          // 日志配置
      "PrintConsole": true,
      "DetailMode": true,
      "EnablePushLog": true
    }
  },
  "AsyncTask": {                         // 异步任务配置
    "Port": "8009",                      // 异步任务服务端口
    "Mysql": "DSN连接字符串"
  },
  "CronJob": {                           // 定时任务配置
    "Mysql": "DSN连接字符串"
  }
}
```

## 部署指南

### Docker 部署

项目提供完整的 Docker 部署方案，使用 Supervisord 管理多个服务。

1. **构建 Docker 镜像**
   ```bash
   cd gaia-server
   docker build -f deploy/Dockerfile -t gaia-server:latest .
   ```

2. **运行容器**
   ```bash
   docker run -d \
     --name gaia-server \
     --restart unless-stopped \
     -p 8008:8008 \
     -v $(pwd)/configs/local:/app/configs/local \
     -v $(pwd)/var/logs:/app/var/logs \
     gaia-server:latest
   ```

3. **使用部署脚本**
   ```bash
   ./deploy.sh
   ```

### Supervisord 配置

Supervisord 用于管理容器内的多个服务进程：

```ini
[program:http_server]
command=/app/cmd/run_http_server.sh
autostart=true
autorestart=true
stdout_logfile=/var/log/supervisor/http_server.out.log
stderr_logfile=/var/log/supervisor/http_server.err.log
```

### 环境变量

支持通过环境变量覆盖配置：

```bash
# 设置服务端口
export SERVER_PORT=8080

# 设置数据库连接
export FRAMEWORK_MYSQL="user:pass@tcp(host:port)/db"
```

## 开发指南

### 添加新的 API

1. **创建控制器** 在 `app/server/api/` 目录下添加新的 Go 文件
2. **定义请求结构** 使用 Gaia SDK 的数据检查标签
3. **注册路由** 在 `app/server/router/` 中添加路由注册

示例：
```go
type UserCtrl struct{}

func (c *UserCtrl) Create() app.HandlerFunc {
    return server.MakeHandler(func(arg server.Request) (any, error) {
        req := CreateUserRequest{}
        if err := arg.BindJsonWithChecker(&req); err != nil {
            return nil, err
        }
        // 业务逻辑
        return map[string]any{"id": 1}, nil
    })
}
```

### 添加定时任务

1. **创建任务** 在 `app/server/jobs/` 目录下添加任务定义
2. **注册任务** 在任务文件中使用 `gaia.RegisterCronJob` 注册

### 数据库操作

使用 GORM Gen 生成的 DAO：
```go
import "gaia-server/app/server/repo/dao"

user := dao.TUser
result, err := user.WithContext(ctx).Where(user.ID.Eq(1)).First()
```

### 使用 Gaia SDK 模块

Gaia SDK 提供了丰富的内置模块：

- **配置管理**: `gaia.GetConfig()`
- **缓存操作**: `gaia.GetCache()`
- **数据库**: `gaia.NewFrameworkMysql()`
- **日志**: `gaia.InfoLog()`, `gaia.ErrorLog()`
- **HTTP 客户端**: `httpclient.NewHttpRequest()`
- **数据验证**: 通过 struct tag `require:"1"` 等

## 常见问题

### 1. 如何修改服务端口？
修改 `configs/local/config.json` 中的 `Server.Port` 字段。

### 2. 如何启用跨域支持？
将 `Server.Cors.Enable` 设置为 `true`，并配置 `AllowOrigins`。

### 3. 数据库迁移如何操作？
项目使用 GORM 的 AutoMigrate 功能，启动服务时会自动创建表结构。

### 4. 如何查看服务日志？
- 控制台日志：启动时配置 `Logger.PrintConsole: true`
- 文件日志：日志文件位于 `var/logs/` 目录
- Docker 日志：`docker logs -f gaia-server`

### 5. 如何添加新的服务类型？
在 `bin/service.go` 的 `Service` 结构体中添加对应的方法，然后在 `cmd/` 目录下创建启动脚本。

## 许可证

[待补充]

## 联系方式

[待补充]