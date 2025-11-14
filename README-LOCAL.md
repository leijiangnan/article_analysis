# 本地开发启动指南

## 🚀 一键启动

```bash
./start-local.sh
```

这个命令会：
1. 检查端口占用（8080和5173）
2. 直接使用最新源代码启动后端服务（http://localhost:8080）
3. 直接使用最新源代码启动前端服务（http://localhost:5173）
4. 显示服务状态和访问地址

> ✅ 重要：脚本会使用最新源代码启动，不会使用已编译的二进制文件，确保你能看到最新的代码变更效果！

## 🛑 一键停止

```bash
./stop-local.sh
```

## 🔄 一键重启

```bash
./restart-local.sh
```

## 📋 服务地址

- **前端界面**: http://localhost:5173
- **后端API**: http://localhost:8080/api
- **健康检查**: http://localhost:8080/health

## 📁 日志文件

- **后端日志**: `./logs/backend.log`
- **前端日志**: `./logs/frontend.log`
- **进程PID**: `./logs/backend.pid` 和 `./logs/frontend.pid`

## ⚠️ 环境要求

- Go 1.21+
- Node.js 20+
- npm 或 yarn

## 🔧 手动启动（备用方案）

如果需要手动分别启动：

### 后端
```bash
cd backend
go run cmd/main.go
```

### 前端
```bash
cd frontend
npm run dev
```

## 🐛 常见问题

1. **端口被占用**: 脚本会自动检测，关闭占用端口的程序或修改端口配置
2. **依赖缺失**: 脚本会自动安装npm依赖，Go依赖需要手动下载
3. **编译失败**: 检查Go和Node.js版本是否符合要求