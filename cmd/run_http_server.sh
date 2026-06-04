#!/bin/sh
cd "$(dirname "$0")/.." || exit 1

echo "正在编译服务..."
if ! go build -o bin/server ./bin; then
    echo "编译失败，请检查错误信息"
    exit 1
fi

echo "编译成功，启动服务..."
./bin/server -Service=Server -Arg="8008"