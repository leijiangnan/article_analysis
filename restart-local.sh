#!/bin/bash

# 一键重启本地服务脚本

echo "🔄 正在重启文章分析系统..."

# 停止现有服务
if [ -f "stop-local.sh" ]; then
    ./stop-local.sh
else
    echo "停止脚本不存在，尝试直接停止进程..."
    pkill -f "go run cmd/main.go" 2>/dev/null || true
    pkill -f "npm run dev" 2>/dev/null || true
    pkill -f "vite" 2>/dev/null || true
fi

# 清理可能残留的进程
sleep 2

# 清理日志文件（可选）
read -p "是否清理日志文件？(y/n): " clear_logs
if [[ $clear_logs =~ ^[Yy]$ ]]; then
    echo "清理日志文件..."
    rm -f logs/*.log
fi

# 启动服务
if [ -f "start-local.sh" ]; then
    ./start-local.sh
else
    echo "❌ 启动脚本不存在"
    exit 1
fi