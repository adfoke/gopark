FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gopark ./cmd/main.go

# 使用轻量级的alpine镜像
FROM alpine:latest

# 安装CA证书，用于HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder阶段复制编译好的应用
COPY --from=builder /app/gopark .
# 复制配置文件
COPY --from=builder /app/config ./config
# 复制迁移文件
COPY --from=builder /app/internal/migrations ./internal/migrations

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./gopark"]