# AIMiniProgram tRPC-Go Makefile

.PHONY: help proto proto-go proto-ts clean build run

# 默认目标
help:
	@echo "Available targets:"
	@echo "  proto     - Generate protobuf code for all languages"
	@echo "  proto-go  - Generate Go code from protobuf"
	@echo "  proto-ts  - Generate TypeScript code from protobuf"
	@echo "  build     - Build the Go server"
	@echo "  run       - Run the Go server"
	@echo "  clean     - Clean generated files"

# Proto 文件路径
PROTO_DIR := proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/aiminiprogram/*/v1/*.proto)

# 生成所有代码
proto: proto-go proto-ts

# 生成 Go 代码
proto-go:
	@echo "Generating Go code from protobuf..."
	@mkdir -p backend-trpc/internal/pb
	@protoc \
		--go_out=backend-trpc/internal/pb \
		--go_opt=paths=source_relative \
		--go-trpc_out=backend-trpc/internal/pb \
		--go-trpc_opt=paths=source_relative \
		--proto_path=proto \
		$(PROTO_FILES)
	@echo "Go code generated successfully!"

# 生成 TypeScript 代码
proto-ts:
	@echo "Generating TypeScript code from protobuf..."
	@mkdir -p frontend/src/proto
	@protoc \
		--ts_out=frontend/src/proto \
		--ts_opt=long_type_string,client_grpc1,server_grpc1 \
		--proto_path=proto \
		$(PROTO_FILES)
	@echo "TypeScript code generated successfully!"

# 构建 Go 服务
build:
	@echo "Building Go server..."
	@cd backend-trpc && go build -o bin/server cmd/server/main.go
	@echo "Build complete: backend-trpc/bin/server"

# 运行 Go 服务
run: build
	@./backend-trpc/bin/server

# 清理生成文件
clean:
	@echo "Cleaning generated files..."
	@rm -rf backend-trpc/internal/pb
	@rm -rf frontend/src/proto
	@rm -rf backend-trpc/bin
	@echo "Clean complete!"

# 安装依赖
deps:
	@echo "Installing Go dependencies..."
	@cd backend-trpc && go mod tidy
	@echo "Installing protobuf tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install trpc.group/trpc-go/trpc-cmdline/protoc-gen-go-trpc@latest
	@echo "Dependencies installed!"

# 检查工具
check:
	@echo "Checking required tools..."
	@which protoc || echo "ERROR: protoc not found"
	@which protoc-gen-go || echo "ERROR: protoc-gen-go not found"
	@which protoc-gen-go-trpc || echo "ERROR: protoc-gen-go-trpc not found"
