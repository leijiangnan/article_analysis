#!/bin/bash

# 一键停止本地服务脚本

echo "🛑 正在停止文章分析系统..."

# 停止后端服务
if [ -f "logs/backend.pid" ]; then
    BACKEND_PID=$(cat logs/backend.pid)
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        echo "停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        rm -f logs/backend.pid
    fi
fi

# 停止前端服务
if [ -f "logs/frontend.pid" ]; then
    FRONTEND_PID=$(cat logs/frontend.pid)
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        echo "停止前端服务 (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        rm -f logs/frontend.pid
    fi
fi

# 清理可能残留的进程
echo "清理残留进程..."
pkill -f "go run cmd/main.go" 2>/dev/null || true
pkill -f "npm run dev" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true
pkill -f "node" 2>/dev/null || true

echo "✅ 服务已停止"
echo "📋 日志文件保留在 ./logs/ 目录中"