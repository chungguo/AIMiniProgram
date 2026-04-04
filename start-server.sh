#!/bin/bash

# 启动 tRPC-Go 服务

echo "=== Starting AIMiniProgram tRPC-Go Server ==="

# 检查可执行文件
if [ ! -f "backend-trpc/bin/server" ]; then
    echo "Building server..."
    cd backend-trpc && go build -o bin/server cmd/server/main.go
    cd ..
fi

# 启动服务
echo "Starting server on port 8000..."
./backend-trpc/bin/server
