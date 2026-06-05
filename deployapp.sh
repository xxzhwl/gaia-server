#!/bin/bash
set -e

# ============================================================
# Gaia Server 一键部署脚本
# 功能：启动基础设施(Redis/Prometheus/Grafana) + 构建运行 gaia-server
# ============================================================

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
INFRA_COMPOSE="$PROJECT_ROOT/deploy/docker-compose-infra.yml"
DOCKERFILE="$PROJECT_ROOT/deploy/Dockerfile"
IMAGE_NAME="gaia-server:latest"
CONTAINER_NAME="gaia-server"
NETWORK_NAME="gaia-network"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# ============================================================
# 1. 前置检查
# ============================================================
check_prerequisites() {
    log_info "检查环境依赖..."

    if ! command -v docker &>/dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if docker compose version &>/dev/null; then
        COMPOSE_CMD="docker compose"
    elif command -v docker-compose &>/dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    log_info "使用 Compose 命令: $COMPOSE_CMD"
}

# ============================================================
# 2. 启动基础设施服务
# ============================================================
start_infra() {
    log_info "启动基础设施服务 (Redis, Prometheus, Grafana)..."
    $COMPOSE_CMD -f "$INFRA_COMPOSE" up -d

    log_info "等待 Redis 就绪..."
    for i in {1..20}; do
        if docker exec gaia-redis redis-cli ping &>/dev/null; then
            log_info "Redis 已就绪"
            break
        fi
        if [ $i -eq 20 ]; then
            log_warn "Redis 启动超时，继续部署..."
        fi
        sleep 2
    done

    log_info "等待 Prometheus 就绪..."
    for i in {1..15}; do
        if curl -sf http://localhost:9090/-/ready &>/dev/null; then
            log_info "Prometheus 已就绪"
            break
        fi
        if [ $i -eq 15 ]; then
            log_warn "Prometheus 启动超时，继续部署..."
        fi
        sleep 2
    done

    log_info "等待 Grafana 就绪..."
    for i in {1..15}; do
        if curl -sf http://localhost:3000/api/health &>/dev/null; then
            log_info "Grafana 已就绪"
            break
        fi
        if [ $i -eq 15 ]; then
            log_warn "Grafana 启动超时，继续部署..."
        fi
        sleep 2
    done
}

# ============================================================
# 3. 构建 gaia-server 镜像
# ============================================================
build_image() {
    log_info "构建 gaia-server 镜像..."
    docker build -f "$DOCKERFILE" -t "$IMAGE_NAME" "$PROJECT_ROOT"
    log_info "镜像构建完成: $IMAGE_NAME"
}

# ============================================================
# 4. 运行 gaia-server 容器
# ============================================================
run_container() {
    # 检查并移除已有容器
    if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_info "停止并移除旧容器: $CONTAINER_NAME"
        docker stop "$CONTAINER_NAME" 2>/dev/null || true
        docker rm "$CONTAINER_NAME" 2>/dev/null || true
    fi

    log_info "启动 gaia-server 容器..."
    docker run -d \
        --name "$CONTAINER_NAME" \
        --restart unless-stopped \
        -p 8008:8008 \
        -p 9100:9100 \
        --network "$NETWORK_NAME" \
        "$IMAGE_NAME"

    log_info "等待 gaia-server 启动..."
    sleep 5
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_info "gaia-server 容器已启动"
    else
        log_error "gaia-server 容器启动失败，请检查日志: docker logs $CONTAINER_NAME"
        exit 1
    fi
}

# ============================================================
# 5. 输出部署结果
# ============================================================
print_summary() {
    echo ""
    echo "============================================"
    log_info "部署完成！"
    echo "============================================"
    echo ""
    echo "  gaia-server:  http://localhost:8008"
    echo "  Metrics:      http://localhost:9100/metrics"
    echo "  Prometheus:   http://localhost:9090"
    echo "  Grafana:      http://localhost:3000  (admin/admin)"
    echo ""
    echo "  查看日志:     docker logs -f $CONTAINER_NAME"
    echo "  停止服务:     docker stop $CONTAINER_NAME"
    echo "  停止全部:     $COMPOSE_CMD -f $INFRA_COMPOSE down && docker stop $CONTAINER_NAME"
    echo ""
}

# ============================================================
# 主流程
# ============================================================
main() {
    echo ""
    echo "============================================"
    echo "       Gaia Server 一键部署"
    echo "============================================"
    echo ""

    check_prerequisites
    start_infra
    build_image
    run_container
    print_summary
}

main "$@"
