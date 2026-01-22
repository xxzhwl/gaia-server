# 开发运行阶段 - 使用 supervisord 管理多个服务
FROM golang:alpine
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 安装 supervisord
RUN apk add --no-cache supervisor

WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 复制配置文件到正确位置
COPY configs/local/config.json ./configs/local/config.json

# 复制 supervisord 配置
COPY deploy/supervisord.conf /etc/supervisor/conf.d/goback.conf

# 创建日志目录
RUN mkdir -p /var/log/supervisor

# 调试：显示 cmd 目录内容
RUN ls -la /app/cmd/

# 确保 shell 脚本有执行权限
RUN chmod +x /app/cmd/*.sh

# 暴露端口
EXPOSE 8008

# 启动 supervisord
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/gaia-server.conf"]