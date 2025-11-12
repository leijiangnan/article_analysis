# Docker部署指南

本文档介绍如何使用Docker快速部署文章分析系统。

## 🚀 快速开始

### 前提条件
- Docker 20.10+
- Docker Compose 1.29+

### 一键启动
```bash
# 克隆项目（如果尚未克隆）
git clone <your-repo-url>
cd article_analysis

# 启动系统
./start.sh
```

访问地址：
- 前端应用：http://localhost
- 后端API：http://localhost:8080/api
- 健康检查：http://localhost:8080/api/health

## 📋 手动部署

### 1. 环境准备
```bash
# 创建环境变量文件
cp .env.example .env

# 根据需要修改配置
nano .env
```

### 2. 构建镜像
```bash
# 构建所有镜像
docker-compose build

# 或者分别构建
docker build -t article-analysis-backend ./backend
docker build -t article-analysis-frontend ./frontend
```

### 3. 启动服务
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 4. 停止服务
```bash
# 停止服务
./stop.sh

# 或者手动停止
docker-compose down

# 停止并删除镜像（清理空间）
docker-compose down --rmi all
```

## 🔧 配置说明

### 环境变量
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| BACKEND_PORT | 后端服务端口 | 8080 |
| FRONTEND_PORT | 前端服务端口 | 80 |
| VITE_API_BASE_URL | API基础URL | http://localhost:8080/api |
| GIN_MODE | Gin框架模式 | release |
| LOG_LEVEL | 日志级别 | info |

### 数据持久化
- 后端数据：`backend_data` 卷，包含SQLite数据库文件
- 配置文件：`config.yaml` 挂载为只读

### 网络配置
- 使用自定义网络 `app-network`，子网 `172.20.0.0/16`
- 服务间可通过服务名直接访问（如 `backend:8080`）

## 🏗️ 架构说明

### 服务架构
```
┌─────────────────┐    ┌─────────────────┐
│   Nginx (80)    │    │  Backend (8080) │
│  前端静态文件    │◄───┤   API服务       │
│  + API代理      │    │  + SQLite数据库  │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────┬───────────────┘
                 │
         ┌───────▼────────┐
         │   Docker网络   │
         │  app-network   │
         └────────────────┘
```

### 技术栈
- **前端**：Vue 3 + TypeScript + Vite + Element Plus
- **后端**：Go + Gin + GORM + SQLite
- **Web服务器**：Nginx（前端静态文件 + API代理）
- **容器化**：Docker + Docker Compose

## 🔍 故障排查

### 查看日志
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend

# 实时查看日志
docker-compose logs -f
```

### 进入容器
```bash
# 进入后端容器
docker-compose exec backend sh

# 进入前端容器
docker-compose exec frontend sh
```

### 常见问题

#### 1. 端口被占用
```bash
# 检查端口占用
lsof -i :80
lsof -i :8080

# 修改端口映射
# 编辑 docker-compose.yml 中的 ports 部分
```

#### 2. 构建失败
```bash
# 清理缓存并重新构建
docker-compose build --no-cache

# 检查网络连接
# 确保可以访问Docker Hub和npm源
```

#### 3. 数据库问题
```bash
# 检查数据卷
docker volume ls
docker volume inspect article_analysis_backend_data

# 重置数据库（谨慎操作）
docker-compose down -v  # 删除数据卷
docker-compose up -d    # 重新启动
```

## 🔒 安全建议

### 生产环境配置
1. **修改默认密码**：更新 `.env` 文件中的敏感信息
2. **使用HTTPS**：配置SSL证书
3. **限制访问**：配置防火墙规则
4. **定期备份**：备份数据库和配置文件

### 性能优化
1. **资源限制**：在docker-compose.yml中设置内存和CPU限制
2. **日志轮转**：配置日志大小限制
3. **监控**：集成Prometheus和Grafana监控

## 📊 监控和运维

### 资源使用
```bash
# 查看资源使用
docker stats

# 查看容器信息
docker-compose ps

# 查看镜像大小
docker images
```

### 备份策略
```bash
# 备份数据库
docker-compose exec backend cp /root/data/article_analysis.db /tmp/
docker cp article_analysis_backend:/tmp/article_analysis.db ./backup/

# 备份配置
cp -r ./backend/config.yaml ./backup/
```

## 🆘 获取帮助

如有问题，请：
1. 查看容器日志：`docker-compose logs`
2. 检查服务状态：`docker-compose ps`
3. 提交Issue到项目仓库

## 📄 许可证

本项目采用MIT许可证，详见项目根目录的LICENSE文件。