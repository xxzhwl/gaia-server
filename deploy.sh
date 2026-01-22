#!/bin/bash
set -e

echo "调用部署目录中的脚本..."
exec ./deploy/deploy.sh "$@"