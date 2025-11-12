# 构建阶段
FROM node:20 AS builder

# 设置工作目录
WORKDIR /app

# 复制package文件
COPY package*.json ./

# 安装依赖（包含开发依赖用于构建）
RUN npm ci

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 生产阶段
FROM nginx

# 复制构建产物到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置文件
COPY nginx.conf /etc/nginx/nginx.conf

# 创建日志目录并设置权限
RUN mkdir -p /var/log/nginx /var/cache/nginx && \
    chown -R www-data:www-data /var/log/nginx /var/cache/nginx

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]