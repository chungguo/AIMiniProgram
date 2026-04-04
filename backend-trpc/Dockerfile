# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装 git（用于获取依赖）
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY backend-trpc/go.mod backend-trpc/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY backend-trpc/ .

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates（用于 HTTPS）
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/bin/server .

# 复制数据文件（如果有）
COPY --from=builder /app/internal/repository/data ./data

# 暴露端口
EXPOSE 8000

# 运行服务
CMD ["./server"]
