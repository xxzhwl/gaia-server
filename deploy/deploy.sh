#!/bin/bash
set -e

echo "开始部署 Gaia 服务..."
echo "使用部署目录中的配置..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装。请先安装 Docker。"
    exit 1
fi

# 构建 Docker 镜像（从项目根目录构建，使用指定 Dockerfile）
echo "构建 Docker 镜像..."
cd "$(dirname "$0")/.."
docker build -f deploy/Dockerfile -t gaia-server:latest .

# 停止并删除现有容器（如果存在）
if docker ps -a --format '{{.Names}}' | grep -q '^gaia-server$'; then
    echo "停止并删除现有容器..."
    docker stop gaia-server || true
    docker rm gaia-server || true
fi

# 运行新容器
echo "启动新容器..."
docker run -d \
    --name gaia-server \
    --restart unless-stopped \
    -p 8008:8008 \
    gaia-server:latest

echo "部署完成！"
echo "HTTP 服务运行在端口 8008"
echo "使用以下命令查看日志: docker logs -f gaia-server"