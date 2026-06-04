#!/bin/bash
set -e

echo "=== 开始部署 Gaia 基础设施（ETCD 配置中心）==="

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装。请先安装 Docker。"
    exit 1
fi

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 使用 ETCD 版配置
if [ -f "configs/local/config-etcd.json" ]; then
    cp configs/local/config-etcd.json configs/local/config.json
    echo "已切换配置为 ETCD 版本"
else
    echo "错误: 找不到 configs/local/config-etcd.json"
    exit 1
fi

# 启动依赖服务（ETCD 版 docker-compose）
echo "启动依赖服务 (MySQL, ETCD, Redis, Elasticsearch, Jaeger, Prometheus, Grafana)..."
docker compose -f docker-compose-etcd.yml up -d

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

# 等待 ETCD 就绪
echo "等待 ETCD 就绪..."
for i in $(seq 1 30); do
    if curl -sf http://127.0.0.1:2379/health &>/dev/null; then
        echo "ETCD 已就绪"
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "错误: ETCD 启动超时"
        exit 1
    fi
    sleep 2
done

# 向 ETCD 写入初始配置
# 注意：地址使用 127.0.0.1，因为 gaia-server 在宿主机运行，通过 docker 端口映射访问服务
echo "写入初始配置到 ETCD..."
ETCDCTL="docker exec gaia-etcd etcdctl"

$ETCDCTL put /gaia/app/Framework.Mysql \
  '"root:123456@tcp(127.0.0.1:3306)/account_system?charset=utf8mb4&parseTime=True&loc=Local"'
$ETCDCTL put /gaia/app/Framework.Redis.Address '"127.0.0.1:6379"'
$ETCDCTL put /gaia/app/Framework.Redis.UserName '""'
$ETCDCTL put /gaia/app/Framework.Redis.Password '""'
$ETCDCTL put /gaia/app/Framework.ES.Address '"http://127.0.0.1:9200"'
$ETCDCTL put /gaia/app/Framework.ES.UserName '""'
$ETCDCTL put /gaia/app/Framework.ES.Password '""'
$ETCDCTL put /gaia/app/Framework.JaegerTracePoint '"127.0.0.1:4318"'
$ETCDCTL put /gaia/app/Framework.Metrics.Enabled '"true"'
$ETCDCTL put /gaia/app/Framework.Metrics.Backend '"prometheus"'

$ETCDCTL put /gaia/app/JWT.SecretKey '"your-secret-key-change-this-in-production"'
$ETCDCTL put /gaia/app/JWT.AccessTokenExp '"15"'
$ETCDCTL put /gaia/app/JWT.RefreshTokenExp '"168"'

echo "ETCD 初始配置写入完成"

echo ""
echo "=== 基础设施部署完成 ==="
echo ""
echo "启动 gaia-server（在宿主机运行）："
echo "  go run cmd/main.go"
echo "  或: ./cmd/run_http_server.sh"
echo ""
echo "ETCD 端点:   http://127.0.0.1:2379"
echo "Jaeger UI:   http://127.0.0.1:16686"
echo "Prometheus:  http://localhost:9090"
echo "Grafana:     http://localhost:3000 (admin/admin)"
echo ""
echo "查看 ETCD 配置: docker exec gaia-etcd etcdctl get /gaia/app --prefix"
