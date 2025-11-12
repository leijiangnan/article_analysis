#!/bin/bash

set -e

echo "========== 开始构建项目 =========="

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go，请先安装Go 1.21+"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: 未找到Node.js，请先安装Node.js 20+"
    exit 1
fi

echo "========== 构建后端应用 =========="
cd "$(dirname "$0")/backend"

# 下载Go依赖
echo "正在下载Go依赖..."
go mod download
go mod tidy

# 构建Go应用（启用CGO以支持SQLite）
echo "正在构建Go应用..."
# 使用Docker构建支持CGO的Linux二进制文件
echo "使用Docker构建支持CGO的Linux二进制文件..."
docker run --rm -v "$(pwd):/go/src/app" -w /go/src/app golang:1.21-alpine sh -c "
  apk add --no-cache gcc musl-dev sqlite-dev &&
  go mod download &&
  CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/main.go
"

# 同时构建MacOS版本
echo "构建MacOS版本..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o main-mac ./cmd/main.go

# 检查构建是否成功
if [ ! -f "main" ]; then
    echo "错误: 后端构建失败，未生成main可执行文件"
    exit 1
fi

echo "后端构建成功！"

echo "========== 构建前端应用 =========="
# 切换到前端目录
cd "/Users/bytedance/goproject/my/article_analysis/frontend" || {
  echo "无法进入前端目录！"
  exit 1
}

# 安装npm依赖
echo "正在安装npm依赖..."
npm ci

# 构建前端应用
echo "正在构建前端应用..."
npm run build

# 检查构建是否成功
if [ ! -d "dist" ]; then
    echo "错误: 前端构建失败，未生成dist目录"
    exit 1
fi

echo "前端构建成功！"

echo "========== 构建完成 =========="
echo "所有构建产物已生成："
echo "- 后端二进制文件: ./backend/main"
echo "- 前端构建产物: ./frontend/dist/"
echo ""
echo "现在可以使用 docker-compose up -d 启动服务了"