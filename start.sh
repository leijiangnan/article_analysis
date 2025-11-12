#!/bin/bash

# Docker启动脚本
# 用于快速启动文章分析系统

set -e

echo "🚀 开始启动文章分析系统..."

# 检查Docker和Docker Compose是否已安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 创建环境变量文件（如果不存在）
if [ ! -f .env ]; then
    echo "📋 创建环境变量配置文件..."
    cp .env.example .env
    echo "✅ 已创建.env文件，请根据需要修改配置"
fi

# 构建镜像
echo "🏗️  正在构建Docker镜像..."
docker-compose build

# 启动服务
echo "🚀 正在启动服务..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "🔍 检查服务状态..."
if docker-compose ps | grep -q "Up"; then
    echo "✅ 服务启动成功！"
    echo "📱 前端地址: http://localhost"
    echo "🔧 后端API: http://localhost:8080/api"
    echo "📊 健康检查: http://localhost:8080/api/health"
    echo ""
    echo "📝 常用命令:"
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
    echo "  重启服务: docker-compose restart"
else
    echo "❌ 服务启动失败，请查看日志:"
    docker-compose logs
    exit 1
fi