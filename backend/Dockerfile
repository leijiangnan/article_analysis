# 直接使用alpine作为基础镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs curl

# 设置工作目录
WORKDIR /app

# 从本机构建产物复制二进制文件
COPY main .

# 复制配置文件
COPY config.yaml ./config.yaml

# 复制初始化SQL脚本
COPY scripts/init.sql ./scripts/

# 创建数据目录（在当前工作目录下，避免权限问题）
RUN mkdir -p /app/data

# 设置环境变量，指定配置文件路径
ENV CONFIG_PATH=/app/config.yaml

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# 运行应用
CMD ["./main"]