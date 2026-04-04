#!/bin/bash

# Proto 代码生成脚本

set -e

echo "=== Proto 代码生成 ==="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 protoc
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc is not installed${NC}"
    echo "Please install protoc: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# 安装 Go 插件（如果不存在）
echo "Checking Go plugins..."
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-trpc &> /dev/null; then
    echo "Installing protoc-gen-go-trpc..."
    go install trpc.group/trpc-go/trpc-cmdline/protoc-gen-go-trpc@latest
fi

# 确保输出目录存在
mkdir -p backend-trpc/internal/pb/common/v1
mkdir -p backend-trpc/internal/pb/model/v1
mkdir -p backend-trpc/internal/pb/paper/v1
mkdir -p backend-trpc/internal/pb/analysis/v1

PROTO_DIR="proto"
OUT_DIR="backend-trpc/internal/pb"

echo "Generating Go code from proto files..."

# 生成 common
echo "  - common/v1/common.proto"
protoc \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go-trpc_out=${OUT_DIR} \
    --go-trpc_opt=paths=source_relative \
    --proto_path=${PROTO_DIR} \
    ${PROTO_DIR}/aiminiprogram/common/v1/common.proto

# 生成 model (依赖 common)
echo "  - model/v1/model.proto"
protoc \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go-trpc_out=${OUT_DIR} \
    --go-trpc_opt=paths=source_relative \
    --proto_path=${PROTO_DIR} \
    ${PROTO_DIR}/aiminiprogram/model/v1/model.proto

# 生成 paper (依赖 common)
echo "  - paper/v1/paper.proto"
protoc \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go-trpc_out=${OUT_DIR} \
    --go-trpc_opt=paths=source_relative \
    --proto_path=${PROTO_DIR} \
    ${PROTO_DIR}/aiminiprogram/paper/v1/paper.proto

# 生成 analysis (依赖 common 和 model)
echo "  - analysis/v1/analysis.proto"
protoc \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go-trpc_out=${OUT_DIR} \
    --go-trpc_opt=paths=source_relative \
    --proto_path=${PROTO_DIR} \
    ${PROTO_DIR}/aiminiprogram/analysis/v1/analysis.proto

echo -e "${GREEN}Proto code generated successfully!${NC}"
echo ""
echo "Generated files:"
find ${OUT_DIR} -name "*.pb.go" | head -20
