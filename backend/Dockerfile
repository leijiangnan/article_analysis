# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（启用CGO以支持SQLite）
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs curl

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件到正确位置
COPY --from=builder /app/config.yaml ./config.yaml

# 复制初始化SQL脚本
COPY --from=builder /app/scripts/init.sql ./scripts/

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