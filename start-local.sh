#!/bin/bash

# 一键本地启动脚本
# 同时启动后端和前端服务

set -e

echo "🚀 开始一键启动文章分析系统..."

# 检查端口是否被占用
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null ; then
        echo "❌ 端口 $1 已被占用，请检查其他服务"
        exit 1
    fi
}

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go 1.21+"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "❌ Node.js未安装，请先安装Node.js 20+"
    exit 1
fi

# 检查端口
echo "🔍 检查端口占用..."
check_port 8080
check_port 5173

# 创建日志目录
mkdir -p logs

# 启动后端服务
echo "🏗️  启动后端服务..."
cd backend

# 下载Go依赖（如果需要）
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo "📦 下载Go依赖..."
    go mod download
fi

# 编译后端二进制并启动（避免go run产生的子进程残留）
echo "📦 编译后端二进制..."
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
BIN_NAME="main-local"
echo "目标平台: ${GOOS}/${GOARCH}, 输出二进制: ${BIN_NAME}"
go build -o ${BIN_NAME} ./cmd/main.go
chmod +x ${BIN_NAME}

echo "使用编译后的二进制启动后端..."
nohup ./${BIN_NAME} > ../logs/backend.log 2>&1 &

BACKEND_PID=$!
echo $BACKEND_PID > ../logs/backend.pid

cd ..

# 等待后端启动
echo "⏳ 等待后端服务启动..."
for i in {1..30}; do
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        echo "✅ 后端服务启动成功！"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ 后端服务启动超时，请检查日志: ./logs/backend.log"
        exit 1
    fi
    sleep 1
done

# 启动前端服务
echo "🏗️  启动前端服务..."
cd frontend

# 安装npm依赖（如果需要）
if [ ! -d "node_modules" ]; then
    echo "📦 安装npm依赖..."
    npm install
fi

# 直接使用最新源代码启动前端开发服务器
echo "使用最新源代码启动前端开发服务器..."
nohup npm run dev -- --force > ../logs/frontend.log 2>&1 &

FRONTEND_PID=$!
echo $FRONTEND_PID > ../logs/frontend.pid

cd ..

# 等待前端启动
echo "⏳ 等待前端服务启动..."
sleep 8

# 检查前端是否成功启动
if ! ps -p $FRONTEND_PID > /dev/null; then
    echo "❌ 前端服务启动失败，请检查日志: ./logs/frontend.log"
    exit 1
fi

# 检查服务状态
echo "🔍 检查服务状态..."
echo ""
echo "========================================="
if ps -p $BACKEND_PID > /dev/null; then
    echo "✅ 后端服务运行正常 (PID: $BACKEND_PID)"
    echo "📊 后端API: http://localhost:8080/api"
    echo "🏥 健康检查: http://localhost:8080/health"
else
    echo "❌ 后端服务异常"
fi

echo ""
if ps -p $FRONTEND_PID > /dev/null; then
    echo "✅ 前端服务运行正常 (PID: $FRONTEND_PID)"
    echo "📱 前端界面: http://localhost:5173"
else
    echo "❌ 前端服务异常"
fi

echo "========================================="
echo ""
echo "🎉 文章分析系统启动完成！"
echo ""
echo "📋 日志文件："
echo "   后端日志: ./logs/backend.log"
echo "   前端日志: ./logs/frontend.log"
echo ""
echo "🛑 停止服务："
echo "   执行: ./stop-local.sh"
echo "   或直接按 Ctrl+C"
echo ""
echo "🔄 重启服务："
echo "   执行: ./restart-local.sh"
