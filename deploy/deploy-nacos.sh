#!/bin/bash
set -e

echo "=== 开始部署 Gaia 基础设施（Nacos 配置中心）==="

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装。请先安装 Docker。"
    exit 1
fi

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 使用 Nacos 版配置
if [ -f "configs/local/config-nacos.json" ]; then
    cp configs/local/config-nacos.json configs/local/config.json
    echo "已切换配置为 Nacos 版本"
else
    echo "错误: 找不到 configs/local/config-nacos.json"
    exit 1
fi

# 启动依赖服务（Nacos 版 docker-compose）
echo "启动依赖服务 (MySQL, Nacos, Redis, Elasticsearch, Jaeger, Prometheus, Grafana)..."
docker compose -f docker-compose-nacos.yml up -d

# 等待 MySQL 就绪
echo "等待 MySQL 就绪..."
for i in $(seq 1 30); do
    if docker exec gaia-mysql mysqladmin ping -h 127.0.0.1 -uroot -p123456 &>/dev/null; then
        echo "MySQL 已就绪"
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "错误: MySQL 启动超时"
        exit 1
    fi
    sleep 2
done

# 等待 Nacos 就绪
echo "等待 Nacos 就绪..."
for i in $(seq 1 30); do
    if curl -sf http://127.0.0.1:8848/nacos/v1/console/health/readiness &>/dev/null; then
        echo "Nacos 已就绪"
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "警告: Nacos 启动超时，继续部署..."
    fi
    sleep 3
done

# 向 Nacos 写入初始配置
# 注意：地址使用 127.0.0.1，因为 gaia-server 在宿主机运行，通过 docker 端口映射访问服务
echo "写入初始配置到 Nacos..."
NACOS_CONFIG_CONTENT='Framework:
  Mysql: "root:123456@tcp(127.0.0.1:3306)/account_system?charset=utf8mb4&parseTime=True&loc=Local"
  Redis:
    Address: "127.0.0.1:6379"
    UserName: ""
    Password: ""
  ES:
    Address: "http://127.0.0.1:9200"
    UserName: ""
    Password: ""
  JaegerTracePoint: "127.0.0.1:4318"
  Metrics:
    Enabled: true
    Backend: prometheus
    Prometheus:
      ListenAddr: ":9100"
      Path: "/metrics"
JWT:
  SecretKey: "your-secret-key-change-this-in-production"
  AccessTokenExp: "15"
  RefreshTokenExp: "168"'

curl -sf -X POST "http://127.0.0.1:8848/nacos/v1/cs/configs" \
  -d "dataId=gaia-server.yaml" \
  -d "group=DEFAULT_GROUP" \
  -d "content=$NACOS_CONFIG_CONTENT" \
  &>/dev/null && echo "Nacos 初始配置写入成功" || echo "警告: Nacos 初始配置写入失败（可能已存在，忽略）"

echo ""
echo "=== 基础设施部署完成 ==="
echo ""
echo "启动 gaia-server（在宿主机运行）："
echo "  go run cmd/main.go"
echo "  或: ./cmd/run_http_server.sh"
echo ""
echo "Nacos 控制台: http://127.0.0.1:8848/nacos (用户名/密码: nacos/nacos)"
echo "Jaeger UI:     http://127.0.0.1:16686"
echo "Prometheus:    http://localhost:9090"
echo "Grafana:       http://localhost:3000 (admin/admin)"
echo ""
echo "配置管理："
echo "  - 修改 Nacos 中的 gaia-server.yaml 即可热更新配置"
echo "  - Nacos 控制台 → 配置管理 → 配置列表 → DEFAULT_GROUP → gaia-server.yaml"
