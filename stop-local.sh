#!/bin/bash

# 一键停止本地服务脚本

echo "🛑 正在停止文章分析系统..."

# 停止后端服务
if [ -f "logs/backend.pid" ]; then
    BACKEND_PID=$(cat logs/backend.pid)
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        echo "停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        sleep 2
        # 如果普通kill失败，强制kill
        if ps -p $BACKEND_PID > /dev/null 2>&1; then
            echo "强制停止后端服务..."
            kill -9 $BACKEND_PID
        fi
        rm -f logs/backend.pid
    else
        echo "后端PID文件存在但进程不存在，清理PID文件..."
        rm -f logs/backend.pid
    fi
else
    echo "未找到后端PID文件，尝试通过端口查找进程..."
    # 通过端口查找并停止后端进程
    BACKEND_PID=$(lsof -ti :8080 2>/dev/null)
    if [ ! -z "$BACKEND_PID" ]; then
        echo "找到占用8080端口的进程 (PID: $BACKEND_PID)，正在停止..."
        kill $BACKEND_PID 2>/dev/null || true
        sleep 2
        # 如果普通kill失败，强制kill
        if ps -p $BACKEND_PID > /dev/null 2>&1; then
            echo "强制停止后端进程..."
            kill -9 $BACKEND_PID
        fi
    fi
fi

# 停止前端服务
if [ -f "logs/frontend.pid" ]; then
    FRONTEND_PID=$(cat logs/frontend.pid)
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        echo "停止前端服务 (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        sleep 2
        # 如果普通kill失败，强制kill
        if ps -p $FRONTEND_PID > /dev/null 2>&1; then
            echo "强制停止前端服务..."
            kill -9 $FRONTEND_PID
        fi
        rm -f logs/frontend.pid
    else
        echo "前端PID文件存在但进程不存在，清理PID文件..."
        rm -f logs/frontend.pid
    fi
else
    echo "未找到前端PID文件，尝试通过端口查找进程..."
    # 通过端口查找并停止前端进程
    FRONTEND_PID=$(lsof -ti :5173 2>/dev/null)
    if [ ! -z "$FRONTEND_PID" ]; then
        echo "找到占用5173端口的进程 (PID: $FRONTEND_PID)，正在停止..."
        kill $FRONTEND_PID 2>/dev/null || true
        sleep 2
        # 如果普通kill失败，强制kill
        if ps -p $FRONTEND_PID > /dev/null 2>&1; then
            echo "强制停止前端进程..."
            kill -9 $FRONTEND_PID
        fi
    fi
fi

# 清理可能残留的进程
echo "清理残留进程..."
pkill -f "go run cmd/main.go" 2>/dev/null || true
pkill -f "npm run dev" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true
pkill -f "node" 2>/dev/null || true

# 清理Go编译缓存中的残留进程
GO_CACHE_PIDS=$(ps aux | grep "/Library/Caches/go-build" | grep -v grep | awk '{print $2}')
if [ ! -z "$GO_CACHE_PIDS" ]; then
    echo "清理Go编译缓存中的残留进程..."
    echo "$GO_CACHE_PIDS" | xargs kill -9 2>/dev/null || true
fi

# 等待进程完全退出
echo "等待进程完全退出..."
sleep 3

# 验证端口是否已释放
echo "验证端口状态..."
if lsof -ti :8080 >/dev/null 2>&1; then
    echo "⚠️  警告：8080端口仍被占用"
    echo "剩余进程: $(lsof -ti :8080)"
else
    echo "✅ 8080端口已释放"
fi

if lsof -ti :5173 >/dev/null 2>&1; then
    echo "⚠️  警告：5173端口仍被占用"
    echo "剩余进程: $(lsof -ti :5173)"
else
    echo "✅ 5173端口已释放"
fi

echo "✅ 停止操作完成"
echo "📋 日志文件保留在 ./logs/ 目录中"