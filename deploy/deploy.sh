#!/bin/bash
set -e

echo "开始部署 Gaia 服务..."
echo "使用部署目录中的配置..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装。请先安装 Docker。"
    exit 1
fi

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 启动依赖服务（docker-compose）
echo "启动依赖服务 (MySQL, Consul, Redis, Elasticsearch, Jaeger, Prometheus, Grafana)..."
docker-compose up -d

# 等待 MySQL 就绪
echo "等待 MySQL 就绪..."
for i in {1..30}; do
    if docker exec gaia-mysql mysqladmin ping -h localhost -uroot -p123456 &>/dev/null; then
        echo "MySQL 已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "错误: MySQL 启动超时"
        exit 1
    fi
    sleep 2
done

# 等待 Consul 就绪
echo "等待 Consul 就绪..."
for i in {1..30}; do
    if curl -sf http://localhost:8500/v1/status/leader &>/dev/null; then
        echo "Consul 已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "警告: Consul 启动超时，继续部署..."
    fi
    sleep 2
done

# 构建 Docker 镜像（从项目根目录构建，使用指定 Dockerfile）
echo "构建 Docker 镜像..."
docker build -f deploy/Dockerfile -t gaia-server:latest .

# 停止并删除现有容器（如果存在）
if docker ps -a --format '{{.Names}}' | grep -q '^gaia-server$'; then
    echo "停止并删除现有容器..."
    docker stop gaia-server || true
    docker rm gaia-server || true
fi

# 运行新容器（加入 gaia-network 以便访问 docker-compose 服务）
echo "启动新容器..."
docker run -d \
    --name gaia-server \
    --restart unless-stopped \
    -p 8008:8008 \
    --network gaia-network \
    gaia-server:latest

echo "部署完成！"
echo "HTTP 服务运行在端口 8008"
echo "Prometheus:  http://localhost:9090"
echo "Grafana:     http://localhost:3000 (admin/admin)"
echo "使用以下命令查看日志: docker logs -f gaia-server"