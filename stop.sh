#!/bin/bash

# Docker停止脚本
# 用于停止文章分析系统

set -e

echo "🛑 正在停止文章分析系统..."

# 停止并移除容器
docker-compose down

# 可选：移除镜像（如果需要清理空间）
# echo "🧹 正在清理镜像..."
# docker-compose down --rmi all

echo "✅ 系统已停止"

# 显示状态
echo "📊 当前状态:"
docker-compose ps