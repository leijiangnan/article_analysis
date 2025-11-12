# 直接使用nginx作为基础镜像
FROM nginx

# 复制本地构建产物到nginx目录（会通过docker-compose的volumes挂载或build context提供）
COPY dist /usr/share/nginx/html

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